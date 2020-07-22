package chat

import (
	"time"

	"github.com/rdnply/backend-trainee-assignment/internal/user"
)

type Chat struct {
	ID        int          `json:"id"`
	Name      string       `json:"name"`
	Users     []*user.User `json:"users"`
	CreatedAt time.Time    `json:"created_at"`
}

type Storage interface {
	Add(chatName string, userIDs []int) (int, error)
	Find(name string) (*Chat, error)
	Exists(id int) (bool, error)
}
