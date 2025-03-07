package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"url-shortener-golang/api"
	pb "url-shortener-golang/proto"
	"url-shortener-golang/storage"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	storageType := os.Getenv("STORAGE_TYPE")

	// storageType := flag.String("storage", "memory", "Storage type (memory|postgres)")
	// postgresConn := flag.String("postgres-conn", "", "PostgreSQL connection string")
	// flag.Parse()
	var st storage.Storage
	var err error
	switch storageType {
	case "memory":
		st = storage.NewMemoryStorage()
		log.Println("Using in-memory storage")
	case "postgres":
		postgresConn := generatePostgresConn()
		st, err = storage.NewPostgresStorage(postgresConn)
		if err != nil {
			log.Fatalf("Failed to initialize PostgreSQL storage: %v", err)
		}
		defer st.Close()
		log.Println("Using PostgreSQL storage")
	default:
		log.Fatalf("unsupported storage type: %s", storageType)
	}

	httpHandler := api.NewHandler(st)
	mux := http.NewServeMux()
	httpHandler.RegisterRoutes(mux)

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	grpcServer := grpc.NewServer()
	pb.RegisterURLShortenerServer(grpcServer, api.NewGRPCServer(st))
	reflection.Register(grpcServer)
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen gRPC: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Println("Starting HTTP server on :8080")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	go func() {
		log.Println("Starting gRPC server on :50051")
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("gRPC server error: %v", err)
		}
	}()

	<-ctx.Done()

	log.Println("Shutting down servers...")
	httpServer.Shutdown(context.Background())
	grpcServer.GracefulStop()
}

func generatePostgresConn() string {
	return "postgres://" +
		os.Getenv("POSTGRES_USER") + ":" +
		os.Getenv("POSTGRES_PASSWORD") + "@" +
		os.Getenv("POSTGRES_HOST") + ":" +
		os.Getenv("POSTGRES_PORT") + "/" +
		os.Getenv("POSTGRES_DB") + "?sslmode=disable"
}
