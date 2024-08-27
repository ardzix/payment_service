package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"payment-service/api/proto"
	"payment-service/internal/infrastructure/db"
	"payment-service/internal/infrastructure/paymentgateway"
	"payment-service/internal/infrastructure/repository"
	grpcServer "payment-service/internal/interface/grpc"
	restServer "payment-service/internal/interface/rest"
	"payment-service/internal/usecase"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found. Continuing with environment variables from the environment.")
	}

	// Initialize MongoDB client
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI environment variable is not set")
	}
	mongoClient := db.NewMongoClient(mongoURI)

	// Initialize repository
	paymentRepo := repository.NewMongoPaymentRepository(mongoClient)

	// Initialize payment gateway clients
	stripeClient := paymentgateway.NewStripeClient()
	xenditClient := paymentgateway.NewXenditClient()
	dokuClient := paymentgateway.NewDokuClient()

	// Initialize gRPC client for PaymentConfigService
	grpcAddr := os.Getenv("PAYMENT_CONFIG_SERVICE_HOST")
	if grpcAddr == "" {
		log.Fatal("PAYMENT_CONFIG_SERVICE_HOST environment variable is not set")
	}
	grpcPort := os.Getenv("PAYMENT_CONFIG_SERVICE_PORT")
	if grpcPort == "" {
		log.Fatal("PAYMENT_CONFIG_SERVICE_PORT environment variable is not set")
	}

	grpcConn, err := grpc.Dial(grpcAddr+":"+grpcPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to PaymentConfigService: %v", err)
	}
	defer grpcConn.Close()

	grpcTimeout := os.Getenv("GRPC_TIMEOUT")
	if grpcTimeout == "" {
		log.Fatal("GRPC_TIMEOUT environment variable is not set")
	}

	timeoutDuration, err := time.ParseDuration(grpcTimeout)
	if err != nil {
		log.Fatalf("failed to parse GRPC_TIMEOUT: %v", err)
	}

	paymentConfigClient := paymentgateway.NewPaymentConfigClient(grpcConn, timeoutDuration)

	// Initialize use case
	paymentUseCase := usecase.NewPaymentUseCase(stripeClient, xenditClient, dokuClient, paymentRepo, paymentConfigClient)

	// Initialize gRPC handler
	paymentHandler := grpcServer.NewPaymentHandler(paymentUseCase)

	// Set up gRPC server
	grpcServer := grpc.NewServer()
	proto.RegisterPaymentServiceServer(grpcServer, paymentHandler)

	// Start gRPC server
	go func() {
		lis, err := net.Listen("tcp", ":50056")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		log.Println("gRPC server listening on port 50056")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Set up REST server
	restHandler := restServer.NewPaymentHandler(paymentUseCase)
	router := mux.NewRouter()
	router.HandleFunc("/payments", restHandler.CreatePayment).Methods("POST")

	// Start REST server
	httpServer := &http.Server{
		Addr:         ":8084",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Println("REST server listening on port 8084")
	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatalf("failed to start REST server: %v", err)
	}
}
