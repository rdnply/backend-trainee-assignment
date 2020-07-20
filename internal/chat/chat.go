package chat

import (
	"os/user"
	"time"
)

type Chat struct {
	ID        int
	Name      string
	Users     []*user.User
	CreatedAt *time.Time
}

type Storage interface {
	Add(c *Chat) error
	Get() ([]*Chat, error)
}
