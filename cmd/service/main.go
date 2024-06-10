package main

import (
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"

	"payment-service/api/proto"
	"payment-service/internal/infrastructure/db"
	"payment-service/internal/infrastructure/paymentgateway"
	"payment-service/internal/infrastructure/repository"
	grpcServer "payment-service/internal/interface/grpc"
	"payment-service/internal/usecase"

	"google.golang.org/grpc"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Initialize MongoDB client
	mongoClient := db.NewMongoClient(os.Getenv("MONGO_URI"))

	// Initialize repository
	paymentRepo := repository.NewMongoPaymentRepository(mongoClient)

	// Initialize payment gateway clients
	stripeClient := paymentgateway.NewStripeClient()
	xenditClient := paymentgateway.NewXenditClient()

	// Initialize use case
	paymentUseCase := usecase.NewPaymentUseCase(stripeClient, xenditClient, paymentRepo)

	// Initialize gRPC handler
	paymentHandler := grpcServer.NewPaymentHandler(paymentUseCase)

	// Set up gRPC server
	grpcServer := grpc.NewServer()
	proto.RegisterPaymentServiceServer(grpcServer, paymentHandler)

	// Start gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("gRPC server listening on port 50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
