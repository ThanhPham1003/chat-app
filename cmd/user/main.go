package main

import (
	"log"
	"net"

	"github.com/ThanhPham1003/chat-app/internal/db"
	"github.com/ThanhPham1003/chat-app/internal/user"
	proto "github.com/ThanhPham1003/chat-app/pkg/proto/user"
	"google.golang.org/grpc"
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
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	proto.RegisterUserServiceServer(s, user.NewServer(db, config))
	log.Println("User service running on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
