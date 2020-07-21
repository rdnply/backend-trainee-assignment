package fortest

import "github.com/rdnply/backend-trainee-assignment/internal/user"

var _ user.Storage = &MockUserStorage{}

type MockUserStorage struct {
	Items []*user.User
}

func NewUserStorage() *MockUserStorage {
	return &MockUserStorage{}
}

func (m *MockUserStorage) Add(u *user.User) error {
	m.Items = append(m.Items, u)
	u.ID = len(m.Items)
	return nil
}

func (m *MockUserStorage) Find(username string) (*user.User, error) {
	for _, u := range m.Items {
		if u.Username == username {
			return u, nil
		}
	}

	return &user.User{}, nil
}
