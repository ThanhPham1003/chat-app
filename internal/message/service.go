package message

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/ThanhPham1003/chat-app/internal/db"
	"github.com/ThanhPham1003/chat-app/pkg/proto/message"
	"github.com/golang-jwt/jwt/v4"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type Server struct {
	message.UnimplementedMessageServiceServer
	db     *gorm.DB
	redis  *redis.Client
	config db.Config
}

func NewServer(db *gorm.DB, redis *redis.Client, config db.Config) *Server {
	return &Server{db: db, redis: redis, config: config}
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

func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if info.FullMethod != "/message.MessageService/SendMessage" {
		return handler(ctx, req)
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "Missing metadata")
	}
	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return nil, status.Error(codes.Unauthenticated, "Missing authorization token")
	}
	tokenStr := authHeader[0]
	if len(tokenStr) > 7 && tokenStr[:7] == "Bearer " {
		tokenStr = tokenStr[7:]
	}
	// Access JWTSecret from Server config
	server, ok := info.Server.(*Server)
	if !ok {
		return nil, status.Error(codes.Internal, "Unable to access server config")
	}
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(server.config.JWT.Secret), nil
	})
	if err != nil || !token.Valid {
		return nil, status.Error(codes.Unauthenticated, "Invalid token")
	}
	return handler(ctx, req)
}
