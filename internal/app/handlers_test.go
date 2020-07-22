package app

import (
	"net/http"
	"testing"
	"time"

	"github.com/rdnply/backend-trainee-assignment/internal/chat"
	"github.com/rdnply/backend-trainee-assignment/internal/message"
	"github.com/rdnply/backend-trainee-assignment/internal/test"
	"github.com/rdnply/backend-trainee-assignment/internal/user"
)

func TestAPI(t *testing.T) {
	mockApp := appForTest()
	users := []*user.User{
		{1, "name already exist", time.Now()},
		{2, "name", time.Now()},
		{3, "another name", time.Now()},
	}

	chats := []*chat.Chat{
		{1, "name already exist", users, time.Now()},
		{2, "another chat name", []*user.User{users[0], users[1]}, time.Now()},
	}

	messages := []*message.Message{
		{1, chats[0], users[1], "message", time.Now()},
		{2, chats[1], users[0], "another message", time.Now()},
	}

	mockApp.UserStorage = &test.MockUserStorage{Items: users}
	mockApp.ChatStorage = &test.MockChatStorage{Items: chats}
	mockApp.MessageStorage = &test.MockMessageStorage{Items: messages}

	tests := []test.APITestCase{
		{"add user ok", "POST", "/users/add", `{"username":"unique name"}`, mockApp.addUser, http.StatusCreated, `{"user_id":4}`},
		{"add user incorrect json", "POST", "/users/add", `"username":"unique name"}`, mockApp.addUser, http.StatusBadRequest, "*incorrect json*"},
		{"add user busy name", "POST", "/users/add", `{"username":"name already exist"}`, mockApp.addUser, http.StatusBadRequest, "*already exists*"},

		{"add chat ok", "POST", "/chats/add", `{"name":"unique name", "users":[1,2,3]}`, mockApp.addChat, http.StatusCreated, `{"chat_id":3}`},
		{"add chat incorrect json", "POST", "/chats/add", `"name":"unique name", "users":[1,2,3]}`, mockApp.addChat, http.StatusBadRequest, "*incorrect json*"},
		{"add chat busy name", "POST", "/chats/add", `{"name":"name already exist", "users":[1,2,3]}`, mockApp.addChat, http.StatusBadRequest, "*already exists*"},
		{"add chat user not found", "POST", "/chats/add", `{"name":"unique name", "users":[-111,2,3]}`, mockApp.addChat, http.StatusNotFound, "*not found user*"},

		{"add message ok", "POST", "/messages/add", `{"chat":2, "author":2, "text":"message"}`, mockApp.addMessage, http.StatusCreated, `{"message_id":3}`},
		{"add message incorrect json", "POST", "/messages/add", `"chat":2, "author":2, "text":"message"}`, mockApp.addMessage, http.StatusBadRequest, "*incorrect json*"},
		{"add message chat not found", "POST", "/messages/add", `{"chat":-222, "author":2, "text":"message"}`, mockApp.addMessage, http.StatusNotFound, "*not found chat*"},
		{"add message user not found", "POST", "/messages/add", `{"chat":2, "author":-222, "text":"message"}`, mockApp.addMessage, http.StatusNotFound, "*not found user*"},
	}

	for _, tc := range tests {
		test.Endpoint(t, tc)
	}
}

func appForTest() *App {
	return &App{
		Logger: test.Logger(),
	}
}
