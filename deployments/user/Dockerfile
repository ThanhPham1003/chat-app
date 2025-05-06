# Build stage
FROM golang:1.23.5-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /user ./cmd/user/main.go

# Final stage
FROM alpine:3.21

WORKDIR /app

# Copy binary from builder
COPY --from=builder /user .

# Copy configs
COPY configs/config.yaml ./configs/config.yaml

# Expose ports (gRPC: 50051, HTTP: 8080)
EXPOSE 50051 8080

# Command to run
CMD ["/app/user"]