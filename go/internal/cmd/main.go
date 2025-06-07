package main

import (
	"log"
	"net/http"

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

func main() {
	// Initialize database
	config := &db.Config{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "postgres",
		DBName:   "lumo_db",
		SSLMode:  "disable",
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
