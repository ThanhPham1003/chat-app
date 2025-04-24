package message

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/ThanhPham1003/chat-app/pkg/proto/message"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Server struct {
	message.UnimplementedMessageServiceServer
	db    *gorm.DB
	redis *redis.Client
}

func NewServer(db *gorm.DB, redis *redis.Client) *Server {
	return &Server{db: db, redis: redis}
}

func (s *Server) SendMessage(ctx context.Context, req *message.SendMessageRequest) (*message.SendMessageResponse, error) {
	msg := Message{
		SenderID:   req.SenderId,
		ReceiverID: req.ReceiverId,
		Content:    req.Content,
		Timestamp:  time.Now(),
	}
	if err := s.db.Create(&msg).Error; err != nil {
		return nil, err
	}
	err := s.redis.Publish(ctx, "chat:"+req.ReceiverId, req.Content).Err()
	if err != nil {
		return nil, err
	}
	return &message.SendMessageResponse{MessageId: generateID(), Message: "Message sent"}, nil
}

func (s *Server) StreamMessages(req *message.StreamMessagesRequest, stream message.MessageService_StreamMessagesServer) error {
	sub := s.redis.Subscribe(context.Background(), "chat:"+req.UserId)
	defer sub.Close()
	ch := sub.Channel()
	for msg := range ch {
		if err := stream.Send(&message.Message{
			MessageId:  generateID(),
			SenderId:   "system",
			ReceiverId: req.UserId,
			Content:    msg.Payload,
			Timestamp:  time.Now().Format(time.RFC3339),
		}); err != nil {
			return err
		}
	}
	return nil
}

func generateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
