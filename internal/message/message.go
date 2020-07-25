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
	// Add adds a new message to the storage, updates the message
	// identifier of the transmitted as a parameter.
	Add(m *Message) error

	// GetAll returns a list of all the messages in a particular
	// chat with all the fields sorted by the time the message
	// was created (from early to late).
	GetAll(chatID int) ([]*Message, error)
}
