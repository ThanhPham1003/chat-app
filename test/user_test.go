package test

import (
	"context"
	"testing"

	"github.com/ThanhPham1003/chat-app/internal/db"
	userService "github.com/ThanhPham1003/chat-app/internal/user"
	userProto "github.com/ThanhPham1003/chat-app/pkg/proto/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupUserTest(t *testing.T) (*userService.Server, func()) {
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
	err = db.AutoMigrate(&userService.User{})
	require.NoError(t, err)

	// Create server instance
	server := userService.NewServer(db, config)

	// Cleanup function
	cleanup := func() {
		db.Migrator().DropTable(&userService.User{})
	}

	return server, cleanup
}

func TestRegister(t *testing.T) {
	server, cleanup := setupUserTest(t)
	defer cleanup()

	tests := []struct {
		name           string
		request        *userProto.RegisterRequest
		expectedError  bool
		expectedStatus string
	}{
		{
			name: "successful registration",
			request: &userProto.RegisterRequest{
				Username: "testuser",
				Password: "password123",
			},
			expectedError:  false,
			expectedStatus: "User registered",
		},
		{
			name: "duplicate username",
			request: &userProto.RegisterRequest{
				Username: "testuser",
				Password: "password123",
			},
			expectedError:  false,
			expectedStatus: "Username taken",
		},
		{
			name: "empty username",
			request: &userProto.RegisterRequest{
				Username: "",
				Password: "password123",
			},
			expectedError:  false,
			expectedStatus: "User registered",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := server.Register(context.Background(), tt.request)
			if tt.expectedError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, response.Message)
			if tt.expectedStatus == "User registered" {
				assert.NotEmpty(t, response.UserId)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	server, cleanup := setupUserTest(t)
	defer cleanup()

	// Register a test user first
	_, err := server.Register(context.Background(), &userProto.RegisterRequest{
		Username: "testuser",
		Password: "password123",
	})
	require.NoError(t, err)

	tests := []struct {
		name           string
		request        *userProto.LoginRequest
		expectedError  bool
		expectedStatus string
		expectToken    bool
	}{
		{
			name: "successful login",
			request: &userProto.LoginRequest{
				Username: "testuser",
				Password: "password123",
			},
			expectedError:  false,
			expectedStatus: "Login successful",
			expectToken:    true,
		},
		{
			name: "wrong password",
			request: &userProto.LoginRequest{
				Username: "testuser",
				Password: "wrongpassword",
			},
			expectedError:  false,
			expectedStatus: "Invalid password",
			expectToken:    false,
		},
		{
			name: "non-existent user",
			request: &userProto.LoginRequest{
				Username: "nonexistent",
				Password: "password123",
			},
			expectedError:  false,
			expectedStatus: "User not found",
			expectToken:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := server.Login(context.Background(), tt.request)
			if tt.expectedError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedStatus, response.Message)
			if tt.expectToken {
				assert.NotEmpty(t, response.Token)
			} else {
				assert.Empty(t, response.Token)
			}
		})
	}
}
