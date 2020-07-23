package chat

import "github.com/rdnply/backend-trainee-assignment/internal/format"

type Chat struct {
	ID        int             `json:"id"`
	Name      string          `json:"name"`
	UsersIDs  []int           `json:"users"`
	CreatedAt format.JSONTime `json:"created_at"`
}

type Storage interface {
	Add(c *Chat) error
	Find(name string) (*Chat, error)
	Exists(id int) (bool, error)
	GetAll(id int) ([]*Chat, error)
}
