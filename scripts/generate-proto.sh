  # Create output directories
  mkdir -p pkg/proto/user
  mkdir -p pkg/proto/message

  # Generate user proto files
  protoc \
    --go_out=pkg/proto/user \
    --go_opt=paths=source_relative \
    --go-grpc_out=pkg/proto/user \
    --go-grpc_opt=paths=source_relative \
    --proto_path=api \
    api/user.proto

  # Generate message proto files
  protoc \
    --go_out=pkg/proto/message \
    --go_opt=paths=source_relative \
    --go-grpc_out=pkg/proto/message \
    --go-grpc_opt=paths=source_relative \
    --proto_path=api \
    api/message.proto