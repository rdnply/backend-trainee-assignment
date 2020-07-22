package user

import "time"

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

type Storage interface {
	Add(u *User) error
	Find(username string) (*User, error)
	Exists(id int) (bool, error)
	AllExists(ids []int) (bool, int, error)
}
