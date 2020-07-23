package test

import (
	"github.com/rdnply/backend-trainee-assignment/internal/chat"
	"github.com/rdnply/backend-trainee-assignment/internal/message"
	"github.com/rdnply/backend-trainee-assignment/internal/user"
)

var _ user.Storage = &MockUserStorage{}

type MockUserStorage struct {
	Items []*user.User
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

func (m *MockUserStorage) Exists(id int) (bool, error) {
	for _, u := range m.Items {
		if u.ID == id {
			return true, nil
		}
	}

	return false, nil
}

func (m *MockUserStorage) AllExists(ids []int) (bool, int, error) {
	for _, id := range ids {
		if exists, _ := m.Exists(id); !exists {
			return false, id, nil
		}
	}

	return true, 0, nil
}

var _ chat.Storage = &MockChatStorage{}

type MockChatStorage struct {
	Items []*chat.Chat
}

func (m *MockChatStorage) Add(c *chat.Chat) error {
	m.Items = append(m.Items, c)
	c.ID = len(m.Items)
	return nil
}

func (m *MockChatStorage) Find(name string) (*chat.Chat, error) {
	for _, c := range m.Items {
		if c.Name == name {
			return c, nil
		}
	}

	return &chat.Chat{}, nil
}

func (m *MockChatStorage) Exists(id int) (bool, error) {
	for _, c := range m.Items {
		if c.ID == id {
			return true, nil
		}
	}

	return false, nil
}

func (m *MockChatStorage) GetAll(id int) ([]*chat.Chat, error) {
	chats := make([]*chat.Chat, 0)

	for _, chat := range m.Items {
		for _, userID := range chat.UsersIDs {
			if id == userID {
				chats = append(chats, chat)
				break
			}
		}
	}

	return chats, nil
}

var _ message.Storage = &MockMessageStorage{}

type MockMessageStorage struct {
	Items []*message.Message
}

func (m *MockMessageStorage) Add(msg *message.Message) error {
	m.Items = append(m.Items, msg)
	msg.ID = len(m.Items)
	return nil
}

func (m *MockMessageStorage) GetAll() ([]*message.Message, error) {
	panic("not implemented") // TODO: Implement
}
