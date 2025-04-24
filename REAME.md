# Chat App
A microservices-based chat application using Golang, gRPC, PostgreSQL, and Redis.

## Structure
- `cmd/`: Service binaries (user-service, message-service, client).
- `internal/`: Private service logic and DB utilities.
- `pkg/`: Reusable libraries (proto, auth).
- `api/`: gRPC proto definitions.

## Setup
1. Install dependencies: `go mod download`.
2. Start PostgreSQL/Redis: `docker-compose up -d postgres redis`.
3. Generate proto: `./scripts/generate-proto.sh`.
4. Run services: `go run cmd/user-service/main.go`, etc.

## Usage
- Register/login via client.
- Send messages in `receiver:message` format.