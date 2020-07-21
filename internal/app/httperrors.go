package app

import (
	"net/http"
)

func (app *App) clientError(w http.ResponseWriter, status int, err error) {
	app.Logger.Errorf("%s", err.Error())
	http.Error(w, http.StatusText(status), status)
}

func (app *App) ServerError(w http.ResponseWriter, err error) {
	app.clientError(w, http.StatusInternalServerError, err)
}

func (app *App) BadRequest(w http.ResponseWriter, err error) {
	app.clientError(w, http.StatusBadRequest, err)
}

func (app *App) NotFound(w http.ResponseWriter, err error) {
	app.clientError(w, http.StatusNotFound, err)
}
