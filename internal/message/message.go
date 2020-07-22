package message

import (
	"time"
)

type Message struct {
	ID        int       `json:"id"`
	ChatID    int       `json:"chat"`
	AuthorID  int       `json:"author"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

type Storage interface {
	Add(m *Message) error
	GetAll() ([]*Message, error)
}
