package app

import (
	"net/http"
	"testing"
	"time"

	"github.com/rdnply/backend-trainee-assignment/internal/chat"
	"github.com/rdnply/backend-trainee-assignment/internal/format"
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
		{1, "already exist", []int{1, 2, 3}, format.NewTime(time.Date(2016, time.August, 15, 0, 0, 0, 0, time.UTC))},
		{2, "another chat name", []int{1, 2}, format.NewTime(time.Date(2017, time.August, 15, 0, 0, 0, 0, time.UTC))},
	}

	messages := []*message.Message{
		{1, 1, 2, "message", format.NewTime(time.Date(2016, time.August, 15, 0, 0, 0, 0, time.UTC))},
		{2, 1, 2, "msg", format.NewTime(time.Date(2017, time.August, 15, 0, 0, 0, 0, time.UTC))},
		{3, 2, 1, "another message", format.NewTime(time.Date(2018, time.August, 15, 0, 0, 0, 0, time.UTC))},
	}

	mockApp.UserStorage = &test.MockUserStorage{Items: users}
	mockApp.ChatStorage = &test.MockChatStorage{Items: chats}
	mockApp.MessageStorage = &test.MockMessageStorage{Items: messages}

	tests := []test.APITestCase{
		{"add user ok", "POST", "/users/add", `{"username":"unique name"}`,
			mockApp.addUser, http.StatusCreated, `{"user_id":4}`},
		{"add user incorrect json", "POST", "/users/add", `"username":"unique name"}`,
			mockApp.addUser, http.StatusBadRequest, "*incorrect json*"},
		{"add user busy name", "POST", "/users/add", `{"username":"name already exist"}`,
			mockApp.addUser, http.StatusBadRequest, "*already exists*"},

		{"add chat ok", "POST", "/chats/add", `{"name":"unique name", "users":[2,3]}`,
			mockApp.addChat, http.StatusCreated, `{"chat_id":3}`},
		{"add chat incorrect json", "POST", "/chats/add", `"name":"unique name", "users":[1,2,3]}`,
			mockApp.addChat, http.StatusBadRequest, "*incorrect json*"},
		{"add chat busy name", "POST", "/chats/add", `{"name":"already exist", "users":[1,2,3]}`,
			mockApp.addChat, http.StatusBadRequest, "*already exists*"},
		{"add chat user not found", "POST", "/chats/add", `{"name":"unique name", "users":[-111,2,3]}`,
			mockApp.addChat, http.StatusNotFound, "*not found user*"},

		{"add message ok", "POST", "/messages/add", `{"chat":2, "author":2, "text":"message"}`,
			mockApp.addMessage, http.StatusCreated, `{"message_id":4}`},
		{"add message incorrect json", "POST", "/messages/add", `"chat":2, "author":2, "text":"message"}`,
			mockApp.addMessage, http.StatusBadRequest, "*incorrect json*"},
		{"add message chat not found", "POST", "/messages/add", `{"chat":-222, "author":2, "text":"message"}`,
			mockApp.addMessage, http.StatusNotFound, "*not found chat*"},
		{"add message user not found", "POST", "/messages/add", `{"chat":2, "author":-222, "text":"message"}`,
			mockApp.addMessage, http.StatusNotFound, "*not found user*"},

		{"get chats ok", "POST", "/chats/get", `{"user":1}`,
			mockApp.getChats, http.StatusOK,
			`[{"id":1,"name":"already exist","users":[1,2,3],"created_at":"2016-08-15T00:00:00Z"},` +
				`{"id":2,"name":"another chat name","users":[1,2],"created_at":"2017-08-15T00:00:00Z"}]`},
		{"get chats incorrect json", "POST", "/chats/get", `"user":1}`,
			mockApp.getChats, http.StatusBadRequest, "*incorrect json*"},
		{"get chats user not found", "POST", "/chats/get", `{"user":-111}`,
			mockApp.getChats, http.StatusNotFound, "*not found user*"},

		{"get messages ok", "POST", "/messages/get", `{"chat":1}`,
			mockApp.getMessages, http.StatusOK,
			`[{"id":2,"chat":1,"author":2,"text":"msg","created_at":"2017-08-15T00:00:00Z"},` +
				`{"id":1,"chat":1,"author":2,"text":"message","created_at":"2016-08-15T00:00:00Z"}]`},
		{"get messages incorrect json", "POST", "/messages/get", `"chat":1}`,
			mockApp.getMessages, http.StatusBadRequest, "*incorrect json*"},
		{"get messages chat not found", "POST", "/messages/get", `{"chat":-111}`,
			mockApp.getMessages, http.StatusNotFound, "*not found chat*"},
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
