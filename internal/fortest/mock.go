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
