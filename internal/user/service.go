package user

import (
	"context"
	"crypto/rand"
	"encoding/hex"

	"github.com/ThanhPham1003/chat-app/internal/db"
	"github.com/ThanhPham1003/chat-app/pkg/auth"
	"github.com/ThanhPham1003/chat-app/pkg/proto/user"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Server struct {
	user.UnimplementedUserServiceServer
	db     *gorm.DB
	config *db.Config
}

func NewServer(db *gorm.DB, config *db.Config) *Server {
	return &Server{db: db, config: config}
}

func (s *Server) Register(ctx context.Context, req *user.RegisterRequest) (*user.RegisterResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	u := User{Username: req.Username, Password: string(hashedPassword)}
	if err := s.db.Create(&u).Error; err != nil {
		return &user.RegisterResponse{Message: "Username taken"}, nil
	}
	return &user.RegisterResponse{UserId: generateID(), Message: "User registered"}, nil
}

func (s *Server) Login(ctx context.Context, req *user.LoginRequest) (*user.LoginResponse, error) {
	var u User
	if err := s.db.Where("username = ?", req.Username).First(&u).Error; err != nil {
		return &user.LoginResponse{Message: "User not found"}, nil
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		return &user.LoginResponse{Message: "Invalid password"}, nil
	}
	token, err := auth.GenerateJWT(req.Username, s.config.JWT.Secret)
	if err != nil {
		return nil, err
	}
	return &user.LoginResponse{Token: token, Message: "Login successful"}, nil
}

func generateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
