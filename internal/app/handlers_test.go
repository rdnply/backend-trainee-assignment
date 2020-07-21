package app

import (
	"net/http"
	"testing"

	"github.com/rdnply/backend-trainee-assignment/internal/format"
	"github.com/rdnply/backend-trainee-assignment/internal/fortest"
	"github.com/rdnply/backend-trainee-assignment/internal/user"
)

func TestAPI(t *testing.T) {
	mockApp := appForTest()
	mockApp.UserStorage = &fortest.MockUserStorage{Items: []*user.User{
		{1, "test name", format.NewNullTime()},
		{2, "another name", format.NewNullTime()},
	}}

	tests := []fortest.APITestCase{
		{"add user ok", "POST", "/users/add", `{"username":"unique name"}`, mockApp.addUser, http.StatusOK, `{"id":3}`},
		{"add user incorrect json", "POST", "/users/add", `"username":"unique name"}`, mockApp.addUser, http.StatusBadRequest, ""},
	}

	for _, tc := range tests {
		fortest.Endpoint(t, tc)
	}
}

func appForTest() *App {
	return &App{
		Logger: fortest.Logger(),
	}
}
