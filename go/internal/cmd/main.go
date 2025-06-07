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

	lumeApp "github.com/mcdev12/lumo/go/internal/app/lume"
	"github.com/mcdev12/lumo/go/internal/genproto/protobuf/lume/lumeconnect"
	"github.com/mcdev12/lumo/go/internal/repository/db"
	lumeRepo "github.com/mcdev12/lumo/go/internal/repository/lume"
	lumeService "github.com/mcdev12/lumo/go/internal/service/lume"
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
	repo := lumeRepo.NewRepository(dbConn)
	app := lumeApp.NewLumeApp(repo)
	lumeSvc := lumeService.NewService(app)

	// Create Connect adapter for the Lume service
	lumeServicePath, lumeConnectSvc := lumeconnect.NewLumeServiceHandler(
		lumeSvc,
		connect.WithInterceptors(
		// Add your interceptors here
		),
	)

	// Set up CORS middleware
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

	// Set up routes
	mux := http.NewServeMux()
	mux.Handle(lumeServicePath, lumeConnectSvc)

	// Register gRPC reflection service for grpcui
	reflector := grpcreflect.NewStaticReflector(lumeconnect.LumeServiceName)
	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))

	// Use h2c to support HTTP/2 without TLS
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
