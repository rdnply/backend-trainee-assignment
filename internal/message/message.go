package message

import (
	"time"

	"github.com/rdnply/backend-trainee-assignment/internal/chat"
	"github.com/rdnply/backend-trainee-assignment/internal/user"
)

type Message struct {
	ID        int        `json:"id"`
	Chat      *chat.Chat `json:"chat"`
	Author    *user.User `json:"author"`
	Text      string     `json:"text"`
	CreatedAt time.Time  `json:"created_at"`
}

type Storage interface {
	Add(chatID, authorID int, text string) (int, error)
	GetAll() ([]*Message, error)
}
