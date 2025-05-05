package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/ThanhPham1003/chat-app/internal/db"
	"github.com/ThanhPham1003/chat-app/internal/user"
	userProto "github.com/ThanhPham1003/chat-app/pkg/proto/user"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	config, err := db.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	db, err := db.NewPostgresDB(config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	db.AutoMigrate(&user.User{})

	// Start gRPC server
	grpcLis, err := net.Listen("tcp", config.Services.User.GrpcAddr)
	if err != nil {
		log.Fatalf("Failed to listen (gRPC): %v", err)
	}
	grpcServer := grpc.NewServer()
	userProto.RegisterUserServiceServer(grpcServer, user.NewServer(db, config))

	ctx := context.Background()
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err = userProto.RegisterUserServiceHandlerFromEndpoint(ctx, mux, config.Services.User.GrpcAddr, opts)
	if err != nil {
		log.Fatalf("Failed to register gateway: %v", err)
	}

	// Run servers concurrently
	go func() {
		log.Printf("Message service (gRPC) running on %s", config.Services.User.GrpcAddr)
		if err := grpcServer.Serve(grpcLis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()
	go func() {
		log.Printf("Message service (HTTP) running on %s", config.Services.User.HttpAddr)
		if err := http.ListenAndServe(config.Services.User.HttpAddr, mux); err != nil {
			log.Fatalf("Failed to serve HTTP: %v", err)
		}
	}()

	// Block forever
	select {}
}
