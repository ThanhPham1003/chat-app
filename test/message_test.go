package test

import (
	"context"
	"testing"

	"github.com/ThanhPham1003/chat-app/internal/db"
	msgService "github.com/ThanhPham1003/chat-app/internal/message"
	msgProto "github.com/ThanhPham1003/chat-app/pkg/proto/message"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func init() {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	msgProto.RegisterMessageServiceServer(s, &msgService.Server{})
	go func() {
		if err := s.Serve(lis); err != nil {
			panic(err)
		}
	}()
}

func setupMessageTest(t *testing.T) (*msgService.Server, func()) {
	// Setup test database
	config := &db.Config{}
	config.Database.Host = "localhost"
	config.Database.Port = 5433
	config.Database.User = "phamthanh"
	config.Database.Password = "password"
	config.Database.DBName = "chat_app"
	config.JWT.Secret = "secret-key"

	db, err := db.NewPostgresDB(config)
	require.NoError(t, err)

	// Run migrations
	err = db.AutoMigrate(&msgService.Message{})
	require.NoError(t, err)

	// Setup Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Create server instance
	server := msgService.NewServer(db, redisClient, *config)

	// Cleanup function
	cleanup := func() {
		db.Migrator().DropTable(&msgService.Message{})
		redisClient.Close()
	}

	return server, cleanup
}

func TestSendMessage(t *testing.T) {
	server, cleanup := setupMessageTest(t)
	defer cleanup()

	tests := []struct {
		name           string
		request        *msgProto.SendMessageRequest
		expectedError  bool
		expectedStatus string
	}{
		{
			name: "successful message send",
			request: &msgProto.SendMessageRequest{
				SenderId:   "user1",
				ReceiverId: "user2",
				Content:    "Hello, World!",
			},
			expectedError:  false,
			expectedStatus: "Message sent",
		},
		{
			name: "empty content",
			request: &msgProto.SendMessageRequest{
				SenderId:   "user1",
				ReceiverId: "user2",
				Content:    "",
			},
			expectedError:  false,
			expectedStatus: "Message sent",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := server.SendMessage(context.Background(), tt.request)
			if tt.expectedError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotEmpty(t, response.MessageId)
			assert.Equal(t, tt.expectedStatus, response.Message)
		})
	}
}

// func TestStreamMessages(t *testing.T) {
// 	server, cleanup := setupMessageTest(t)
// 	defer cleanup()

// 	// Create a test stream
// 	stream := &mockMessageStream{
// 		messages: make(chan *msgProto.Message, 1),
// 	}

// 	// Start streaming in a goroutine
// 	go func() {
// 		err := server.StreamMessages(&msgProto.StreamMessagesRequest{
// 			UserId: "user1",
// 		}, stream)
// 		assert.NoError(t, err)
// 	}()

// 	// Send a test message
// 	_, err := server.SendMessage(context.Background(), &msgProto.SendMessageRequest{
// 		SenderId:   "user2",
// 		ReceiverId: "user1",
// 		Content:    "Test message",
// 	})
// 	require.NoError(t, err)

// 	// Wait for the message to be received
// 	select {
// 	case msg := <-stream.messages:
// 		assert.Equal(t, "Test message", msg.Content)
// 		assert.Equal(t, "user2", msg.SenderId)
// 		assert.Equal(t, "user1", msg.ReceiverId)
// 	case <-time.After(5 * time.Second):
// 		t.Fatal("Timeout waiting for message")
// 	}
// }

// // mockMessageStream implements the MessageService_StreamMessagesServer interface
// type mockMessageStream struct {
// 	messages chan *msgProto.Message
// }

// func (m *mockMessageStream) Send(msg *msgProto.Message) error {
// 	m.messages <- msg
// 	return nil
// }

// func (m *mockMessageStream) SetHeader(metadata.MD) error  { return nil }
// func (m *mockMessageStream) SendHeader(metadata.MD) error { return nil }
// func (m *mockMessageStream) SetTrailer(metadata.MD)       {}
// func (m *mockMessageStream) Context() context.Context     { return context.Background() }
// func (m *mockMessageStream) SendMsg(interface{}) error    { return nil }
// func (m *mockMessageStream) RecvMsg(interface{}) error    { return nil }
