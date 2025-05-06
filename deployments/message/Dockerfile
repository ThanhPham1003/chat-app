# Build stage
FROM golang:1.23.5-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /message ./cmd/message/main.go

# Final stage
FROM alpine:3.21

WORKDIR /app

# Copy binary from builder
COPY --from=builder /message .

# Copy configs
COPY configs/config.yaml ./configs/config.yaml

# Expose ports (gRPC: 50052, HTTP: 8081)
EXPOSE 50052 8081

# Command to run
CMD ["/app/message"]