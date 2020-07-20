package app

import (
	"log"
	"net/http"
)

func (app *App) RunServer() {
	srv := &http.Server{
		Addr:    app.addr,
		Handler: app.routes(),
	}

	app.logger.Infof("Server is running on %s", app.addr)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
