package user

import "time"

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

type Storage interface {
	// Add adds a new user to the storage, updates the user
	// identifier of the transmitted as a parameter.
	Add(u *User) error

	// Find returns the user with the passed name, if the
	// user does not exist returns an empty user structure.
	Find(username string) (*User, error)

	// Exists verifies the existence of a user with the passed ID.
	Exists(id int) (bool, error)

	// AllExists checks for the existence of all users in the
	// transmitted list of identifiers, if at least one does
	// not exist then returns false.
	AllExists(ids []int) (bool, int, error)
}
