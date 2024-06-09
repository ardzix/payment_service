package main

import (
	"log"
	"net"
	"os"

	"payment-service/api/proto"
	"payment-service/internal/infrastructure/db"
	"payment-service/internal/infrastructure/paymentgateway"
	"payment-service/internal/infrastructure/repository"
	grpcHandler "payment-service/internal/interface/grpc"
	"payment-service/internal/usecase"

	"google.golang.org/grpc"
)

func main() {
	// Load environment variables
	os.Setenv("STRIPE_API_KEY", "your_stripe_api_key")
	os.Setenv("XENDIT_API_KEY", "your_xendit_api_key")

	// Initialize MongoDB client
	mongoClient := db.NewMongoClient("mongodb://localhost:27017")

	// Initialize repository
	paymentRepo := repository.NewMongoPaymentRepository(mongoClient)

	// Initialize payment gateway clients
	stripeClient := paymentgateway.NewStripeClient()
	xenditClient := paymentgateway.NewXenditClient()

	// Initialize use case
	paymentUseCase := usecase.NewPaymentUseCase(paymentRepo, stripeClient, xenditClient)

	// Initialize gRPC handler
	paymentHandler := grpcHandler.NewPaymentHandler(paymentUseCase)

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
