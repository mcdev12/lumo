package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"connectrpc.com/connect"
	"connectrpc.com/grpcreflect"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	linkApp "github.com/mcdev12/lumo/go/internal/app/link"
	lumeApp "github.com/mcdev12/lumo/go/internal/app/lume"
	lumoApp "github.com/mcdev12/lumo/go/internal/app/lumo"
	linkconnect "github.com/mcdev12/lumo/go/internal/genproto/link/v1/linkv1connect"
	lumeconnect "github.com/mcdev12/lumo/go/internal/genproto/lume/v1/lumev1connect"
	lumoconnect "github.com/mcdev12/lumo/go/internal/genproto/lumo/v1/lumov1connect"
	"github.com/mcdev12/lumo/go/internal/repository/db"
	linkRepo "github.com/mcdev12/lumo/go/internal/repository/link"
	lumeRepo "github.com/mcdev12/lumo/go/internal/repository/lume"
	lumoRepo "github.com/mcdev12/lumo/go/internal/repository/lumo"
	linkService "github.com/mcdev12/lumo/go/internal/service/link"
	lumeService "github.com/mcdev12/lumo/go/internal/service/lume"
	lumoService "github.com/mcdev12/lumo/go/internal/service/lumo"
)

// getEnv returns the value of an environment variable or a default value if not set
func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

// getEnvAsInt returns the value of an environment variable as an integer or a default value if not set
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Warning: Environment variable %s is not a valid integer, using default value %d", key, defaultValue)
		return defaultValue
	}
	return value
}

func main() {
	// Initialize database
	config := &db.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnvAsInt("DB_PORT", 5432),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "lumo_db"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	dbConn, err := db.NewConnection(config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	// Initialize layers
	// Lume service
	lumeRepository := lumeRepo.NewRepository(dbConn)
	lumeApplication := lumeApp.NewLumeApp(lumeRepository)
	lumeSvc := lumeService.NewService(lumeApplication)

	// Lumo service
	lumoRepository := lumoRepo.NewRepository(dbConn)
	lumoApplication := lumoApp.NewLumoApp(lumoRepository)
	lumoSvc := lumoService.NewService(lumoApplication)

	// Link service
	linkRepository := linkRepo.NewRepository(dbConn)
	linkApplication := linkApp.NewLinkApp(linkRepository)
	linkSvc := linkService.NewService(linkApplication)

	// Create Connect adapters
	lumeServicePath, lumeConnectSvc := lumeconnect.NewLumeServiceHandler(
		lumeSvc,
		connect.WithInterceptors(
			// Add your interceptors here
		),
	)
	lumoServicePath, lumoConnectSvc := lumoconnect.NewLumoServiceHandler(
		lumoSvc,
		connect.WithInterceptors(
			// Add your interceptors here
		),
	)
	linkServicePath, linkConnectSvc := linkconnect.NewLinkServiceHandler(
		linkSvc,
		connect.WithInterceptors(
			// Add your interceptors here
		),
	)

	// CORS middleware
	corsMiddleware := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Connect-Protocol-Version")
			w.Header().Set("Access-Control-Max-Age", "3600")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			h.ServeHTTP(w, r)
		})
	}

	// Set up HTTP mux and handlers
	mux := http.NewServeMux()
	mux.Handle(lumeServicePath, lumeConnectSvc)
	mux.Handle(lumoServicePath, lumoConnectSvc)
	mux.Handle(linkServicePath, linkConnectSvc)

	// === Reflection for grpcui/grpcurl ===
	reflector := grpcreflect.NewStaticReflector(
		lumeconnect.LumeServiceName,
		lumoconnect.LumoServiceName,
		linkconnect.LinkServiceName,
	)
	// Register both v1 and v1alpha reflection handlers
	pathV1, handlerV1 := grpcreflect.NewHandlerV1(reflector)
	mux.Handle(pathV1, handlerV1)

	pathAlpha, handlerAlpha := grpcreflect.NewHandlerV1Alpha(reflector)
	mux.Handle(pathAlpha, handlerAlpha)

	// Wrap with CORS + h2c (HTTP/2 without TLS)
	handler := corsMiddleware(mux)
	server := &http.Server{
		Addr:    ":8080",
		Handler: h2c.NewHandler(handler, &http2.Server{}),
	}

	log.Printf("Connect server listening on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
