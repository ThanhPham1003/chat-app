package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/ThanhPham1003/chat-app/internal/db"
	"github.com/ThanhPham1003/chat-app/internal/message"
	massageProto "github.com/ThanhPham1003/chat-app/pkg/proto/message"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	config, err := db.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	dbPostgres, err := db.NewPostgresDB(config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	dbPostgres.AutoMigrate(&message.Message{})
	dbRedis, err := db.NewRedisClient(config)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Start gRPC server
	grpcLis, err := net.Listen("tcp", config.Services.Message.GrpcAddr)
	if err != nil {
		log.Fatalf("Failed to listen (gRPC): %v", err)
	}
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(message.AuthInterceptor))
	massageProto.RegisterMessageServiceServer(grpcServer, message.NewServer(dbPostgres, dbRedis, *config))

	ctx := context.Background()
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err = massageProto.RegisterMessageServiceHandlerFromEndpoint(ctx, mux, config.Services.Message.GrpcAddr, opts)
	if err != nil {
		log.Fatalf("Failed to register gateway: %v", err)
	}

	// Run servers concurrently
	go func() {
		log.Printf("Message service (gRPC) running on %s", config.Services.Message.GrpcAddr)
		if err := grpcServer.Serve(grpcLis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()
	go func() {
		log.Printf("Message service (HTTP) running on %s", config.Services.Message.HttpAddr)
		if err := http.ListenAndServe(config.Services.Message.HttpAddr, mux); err != nil {
			log.Fatalf("Failed to serve HTTP: %v", err)
		}
	}()

	// Block forever
	select {}
}
