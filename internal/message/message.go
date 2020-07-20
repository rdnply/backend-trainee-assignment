package message

import (
	"os/user"
	"time"

	"github.com/rdnply/backend-trainee-assignment/internal/chat"
)

type Message struct {
	ID        int
	Chat      *chat.Chat
	Author    *user.User
	Text      string
	CreatedAt *time.Time
}

type Storage interface {
	Add(m *Message) error
	Get() ([]*Message, error)
}
