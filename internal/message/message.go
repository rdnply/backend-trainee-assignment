package message

import "github.com/rdnply/backend-trainee-assignment/internal/format"

type Message struct {
	ID        int             `json:"id"`
	ChatID    int             `json:"chat"`
	AuthorID  int             `json:"author"`
	Text      string          `json:"text"`
	CreatedAt format.JSONTime `json:"created_at"`
}

type Storage interface {
	Add(m *Message) error
	GetAll() ([]*Message, error)
}
