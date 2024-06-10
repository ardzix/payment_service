package main

import (
	"log"
	"net"
	"os"
	"time"

	"github.com/joho/godotenv"

	"payment-service/api/proto"
	"payment-service/internal/infrastructure/db"
	"payment-service/internal/infrastructure/paymentgateway"
	"payment-service/internal/infrastructure/repository"
	grpcServer "payment-service/internal/interface/grpc"
	"payment-service/internal/usecase"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	dokuClient := paymentgateway.NewDokuClient()

	// Initialize gRPC client for PaymentConfigService
	grpcConn, err := grpc.NewClient(os.Getenv("PAYMENT_CONFIG_SERVICE_ADDRESS"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to PaymentConfigService: %v", err)
	}
	defer grpcConn.Close()

	grpcTimeout, err := time.ParseDuration(os.Getenv("GRPC_TIMEOUT"))
	if err != nil {
		log.Fatalf("failed to parse GRPC_TIMEOUT: %v", err)
	}
	paymentConfigClient := paymentgateway.NewPaymentConfigClient(grpcConn, grpcTimeout)

	// Initialize use case
	paymentUseCase := usecase.NewPaymentUseCase(stripeClient, xenditClient, dokuClient, paymentRepo, paymentConfigClient)

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
