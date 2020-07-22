package app

import (
	"net/http"

	"github.com/go-chi/chi"
)

func (app *App) routes() http.Handler {
	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Post("/users/add", app.addUser)
		r.Post("/chats/add", app.addChat)
	})

	return r
}
