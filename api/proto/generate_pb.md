export PATH="$PATH:$(go env GOPATH)/bin"
protoc --go_out=. --go-grpc_out=. api/proto/payment.proto
protoc --go_out=. --go-grpc_out=. api/proto/paymentconfig.proto