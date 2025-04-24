package message

import (
	"time"

	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	SenderID   string
	ReceiverID string
	Content    string
	Timestamp  time.Time
}
