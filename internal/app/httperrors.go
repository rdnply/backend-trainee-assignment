package app

import (
	"net/http"

	"github.com/pkg/errors"
)

func (app *App) clientError(w http.ResponseWriter, status int, err error, msg string) {
	var logErr, showErr string

	switch {
	case err != nil && msg != "":
		logErr = errors.Wrap(err, msg).Error()
		showErr = `{"error": "` + msg + `"}`
	case err != nil:
		logErr = err.Error()
		showErr = http.StatusText(status)
	case msg != "":
		logErr = msg
		showErr = `{"error": "` + msg + `"}`
	}

	app.Logger.Errorf(logErr)
	http.Error(w, showErr, status)
}

func (app *App) ServerError(w http.ResponseWriter, err error, msg string) {
	app.clientError(w, http.StatusInternalServerError, err, msg)
}

func (app *App) BadRequest(w http.ResponseWriter, err error, msg string) {
	app.clientError(w, http.StatusBadRequest, err, msg)
}

func (app *App) NotFound(w http.ResponseWriter, err error, msg string) {
	app.clientError(w, http.StatusNotFound, err, msg)
}

/* func (app *App) clientError(w http.ResponseWriter, status int, err error) { */
//     app.Logger.Errorf("%s", err.Error())
//     http.Error(w, http.StatusText(status), status)
// }

// func (app *App) ServerError(w http.ResponseWriter, err error) {
//     app.clientError(w, http.StatusInternalServerError, err)
// }

// func (app *App) BadRequest(w http.ResponseWriter, err error) {
//     app.clientError(w, http.StatusBadRequest, err)
// }

// func (app *App) NotFound(w http.ResponseWriter, err error) {
//     app.clientError(w, http.StatusNotFound, err)
/* } */
