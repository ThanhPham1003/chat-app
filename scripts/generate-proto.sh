# Generate user proto files (gRPC and Gateway)
protoc \
  --go_out=pkg/proto/user \
  --go_opt=paths=source_relative \
  --go-grpc_out=pkg/proto/user \
  --go-grpc_opt=paths=source_relative \
  --grpc-gateway_out=pkg/proto/user \
  --grpc-gateway_opt=paths=source_relative \
  --proto_path=api \
  --proto_path=third_party \
  api/user.proto

# Generate message proto files (gRPC and Gateway)
protoc \
  --go_out=pkg/proto/message \
  --go_opt=paths=source_relative \
  --go-grpc_out=pkg/proto/message \
  --go-grpc_opt=paths=source_relative \
  --grpc-gateway_out=pkg/proto/message \
  --grpc-gateway_opt=paths=source_relative \
  --proto_path=api \
  --proto_path=third_party \
  api/message.proto