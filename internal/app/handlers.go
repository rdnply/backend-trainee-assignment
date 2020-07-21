package app

import (
	"encoding/json"
	"net/http"

	"github.com/rdnply/backend-trainee-assignment/internal/user"
)

func (app *App) addUser(w http.ResponseWriter, r *http.Request) {
	var u user.User

	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		app.BadRequest(w, err)
		return
	}

	err = app.UserStorage.Add(&u)
	if err != nil {
		app.ServerError(w, err)
		return
	}

	respondJSON(w, map[string]int{"id": u.ID})
}

func respondJSON(w http.ResponseWriter, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
