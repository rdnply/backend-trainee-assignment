package user

import "github.com/rdnply/backend-trainee-assignment/internal/format"

type User struct {
	ID        int              `json:"id"`
	Username  string           `json:"username"`
	CreatedAt *format.NullTime `json:"created_at"`
}

type Storage interface {
	Add(u *User) error
	// Find(username string) (*User, error)
}
