package main

import (
	"go-lang-final/internal/handlers"
	"go-lang-final/internal/store"
	"log"
	"net"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	logger := logrus.New()
	dsn := "postgresql://user:password@localhost:5432/library"
	paymentStore, err := store.NewPaymentStore(dsn)
	if err != nil {
		logger.Fatalf("Failed to connect to the database: %v", err)
	}

	// REST API
	r := mux.NewRouter()
	handlers.RegisterRESTHandlers(r, paymentStore, logger)

	// gRPC Server
	grpcServer := grpc.NewServer()
	handlers.RegisterGRPCHandlers(grpcServer, paymentStore)

	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			logger.Fatalf("Failed to listen on port 50051: %v", err)
		}
		logger.Info("gRPC server listening on port 50051")
		if err := grpcServer.Serve(lis); err != nil {
			logger.Fatalf("Failed to serve gRPC server: %v", err)
		}
	}()

	logger.Info("Starting HTTP server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
