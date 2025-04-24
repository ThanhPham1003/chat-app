package main

import (
	"log"
	"net"

	"github.com/ThanhPham1003/chat-app/internal/db"
	"github.com/ThanhPham1003/chat-app/internal/message"
	proto "github.com/ThanhPham1003/chat-app/pkg/proto/message"
	"google.golang.org/grpc"
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
	redis, err := db.NewRedisClient(config)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	proto.RegisterMessageServiceServer(s, message.NewServer(dbPostgres, redis))
	log.Println("Message service running on :50052")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
