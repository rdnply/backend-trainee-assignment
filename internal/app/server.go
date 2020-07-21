package app

import (
	"log"
	"net/http"
)

func (app *App) RunServer() {
	srv := &http.Server{
		Addr:    app.Addr,
		Handler: app.routes(),
	}

	app.Logger.Infof("Server is running on %s")
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
