package app

import (
	"net/http"
)

func (app *App) ServerError(w http.ResponseWriter, err error) {
	app.logger.Errorf("%s", err.Error())
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}
