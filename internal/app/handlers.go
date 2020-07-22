package app

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rdnply/backend-trainee-assignment/internal/user"
)

func (app *App) addUser(w http.ResponseWriter, r *http.Request) {
	var u user.User

	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		app.BadRequest(w, err, "incorrect json")
		return
	}

	fromDB, err := app.UserStorage.Find(u.Username)
	if err != nil {
		app.ServerError(w, err, "")
		return
	}
	if fromDB.ID != 0 {
		app.BadRequest(w, err, fmt.Sprintf("user %s already exists", u.Username))
		return
	}

	if err := app.UserStorage.Add(&u); err != nil {
		app.ServerError(w, err, "")
		return
	}

	respondJSON(w, http.StatusCreated, map[string]int{"user_id": u.ID})
}

func (app *App) addChat(w http.ResponseWriter, r *http.Request) {
	type chatInfo struct {
		Name  string `json:"name"`
		Users []int  `json:"users"`
	}

	var c chatInfo
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		app.BadRequest(w, err, "incorrect json")
		return
	}

	exists, id, err := app.UserStorage.AllExists(c.Users)
	if err != nil {
		app.ServerError(w, err, "")
		return
	}
	if !exists {
		app.NotFound(w, err, fmt.Sprintf("not found user with id: %v", id))
		return
	}

	fromDB, err := app.ChatStorage.Find(c.Name)
	if err != nil {
		app.ServerError(w, err, "")
		return
	}
	if fromDB.ID != 0 {
		app.BadRequest(w, err, fmt.Sprintf("chat %s already exists", c.Name))
		return
	}

	newID, err := app.ChatStorage.Add(c.Name, c.Users)
	if err != nil {
		app.ServerError(w, err, "")
		return
	}

	respondJSON(w, http.StatusCreated, map[string]int{"chat_id": newID})
}

func (app *App) addMessage(w http.ResponseWriter, r *http.Request) {
	type msgInfo struct {
		ChatID   int    `json:"chat"`
		AuthorID int    `json:"author"`
		Text     string `json:"text"`
	}

	var m msgInfo
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		app.BadRequest(w, err, "incorrect json")
		return
	}

	existsChat, err := app.ChatStorage.Exists(m.ChatID)
	if err != nil {
		app.ServerError(w, err, "")
		return
	}
	if !existsChat {
		app.NotFound(w, err, fmt.Sprintf("not found chat: %d", m.ChatID))
		return
	}

	existsUser, err := app.UserStorage.Exists(m.AuthorID)
	if err != nil {
		app.ServerError(w, err, "")
		return
	}
	if !existsUser {
		app.NotFound(w, err, fmt.Sprintf("not found user with id: %d", m.AuthorID))
		return
	}

	newID, err := app.MessageStorage.Add(m.ChatID, m.AuthorID, m.Text)
	if err != nil {
		app.ServerError(w, err, "")
		return
	}

	respondJSON(w, http.StatusCreated, map[string]int{"message_id": newID})
}

func respondJSON(w http.ResponseWriter, successCode int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(successCode)
	w.Write(response)
}
