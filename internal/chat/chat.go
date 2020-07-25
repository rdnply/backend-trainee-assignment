package chat

import "github.com/rdnply/backend-trainee-assignment/internal/format"

type Chat struct {
	ID        int             `json:"id"`
	Name      string          `json:"name"`
	UsersIDs  []int           `json:"users"`
	CreatedAt format.JSONTime `json:"created_at"`
}

type Storage interface {
	// Add adds a new chat to the storage, updates the chat
	// identifier of the transmitted as a parameter.
	Add(c *Chat) error

	// Find returns the chat with the passed name, if the
	// chat does not exist returns an empty chat structure.
	Find(name string) (*Chat, error)

	// Exists verifies the existence of a chat with the passed ID.
	Exists(id int) (bool, error)

	// GetAll returns a list of all user chats with the passed ID, sorted
	// by the time the last chat message was created (from late to early).
	GetAll(id int) ([]*Chat, error)
}
