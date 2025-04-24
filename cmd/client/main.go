package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ThanhPham1003/chat-app/pkg/proto/message"
	"github.com/ThanhPham1003/chat-app/pkg/proto/user"
	"google.golang.org/grpc"
)

func main() {
	userConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to user service: %v", err)
	}
	defer userConn.Close()
	userClient := user.NewUserServiceClient(userConn)

	msgConn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to message service: %v", err)
	}
	defer msgConn.Close()
	msgClient := message.NewMessageServiceClient(msgConn)

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)
	fmt.Print("Enter password: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	res, err := userClient.Register(context.Background(), &user.RegisterRequest{Username: username, Password: password})
	if err != nil {
		log.Fatalf("Register failed: %v", err)
	}
	if strings.Contains(res.Message, "taken") {
		_, err = userClient.Login(context.Background(), &user.LoginRequest{Username: username, Password: password})
		if err != nil {
			log.Fatalf("Login failed: %v", err)
		}
	}
	fmt.Println(res.Message)

	stream, err := msgClient.StreamMessages(context.Background(), &message.StreamMessagesRequest{UserId: username})
	if err != nil {
		log.Fatalf("Stream failed: %v", err)
	}
	go func() {
		for {
			msg, err := stream.Recv()
			if err != nil {
				log.Printf("Stream error: %v", err)
				return
			}
			fmt.Printf("[%s] %s: %s\n", msg.Timestamp, msg.SenderId, msg.Content)
		}
	}()

	for {
		fmt.Print("Enter receiver username and message (format: receiver:message): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		parts := strings.SplitN(input, ":", 2)
		if len(parts) != 2 {
			fmt.Println("Invalid format")
			continue
		}
		_, err := msgClient.SendMessage(context.Background(), &message.SendMessageRequest{
			SenderId:   username,
			ReceiverId: parts[0],
			Content:    parts[1],
		})
		if err != nil {
			log.Printf("Send message failed: %v", err)
		}
	}
}
