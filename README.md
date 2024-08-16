
# Payment Service

## Overview

The Payment Service is a gRPC-based service written in Go. It handles various payment operations such as processing payments, refunding payments, retrieving payment status, and listing payments.

## Prerequisites

- Go 1.18 or later
- Docker (optional, for containerized deployment)
- Protocol Buffer Compiler (`protoc`)
- MongoDB

## Setup

### 1. Clone the repository

```sh
git clone <repository-url>
cd payment-service
```

### 2. Install dependencies

```sh
go mod download
```

### 3. Generate gRPC code from protobuf

```sh
protoc --go_out=. --go-grpc_out=. api/proto/payment.proto
protoc --go_out=. --go-grpc_out=. api/proto/paymentconfig.proto
```

### 4. Build and run the service

#### Locally

```sh
go build -o payment-service ./cmd/service
./payment-service
```

#### Using Docker

Build the Docker image:

```sh
docker build -t payment-service .
```

Run the Docker container:

```sh
docker run -p 50056:50056 payment-service
```

## Environment Variables

- `MONGO_URI`: MongoDB connection URI
- `STRIPE_API_KEY`: Stripe API key
- `XENDIT_API_KEY`: Xendit API key
- `PAYMENT_CONFIG_SERVICE_ADDRESS`: Address for payment gateway configuration service
- `GRPC_TIMEOUT`: GRPC timeout in secon
- `DEFAULT_PG`: Default payment gateway if pg configuration service can not be called

## Usage

The service exposes the following gRPC endpoints:

- `ProcessPayment`
- `RefundPayment`
- `GetPaymentStatus`
- `ListPayments`
- `GetPaymentDetail`

Refer to the `payment.proto` file for more details on the request and response formats.

## License

This project is licensed under the MIT License.