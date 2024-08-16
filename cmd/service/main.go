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
	grpcConn, err := grpc.Dial(os.Getenv("PAYMENT_CONFIG_SERVICE_ADDRESS"), grpc.WithTransportCredentials(insecure.NewCredentials()))
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
