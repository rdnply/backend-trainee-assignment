package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rdnply/backend-trainee-assignment/pkg/logger"
)

func (app *App) RunServer() {
	srv := &http.Server{
		Addr:         app.Addr,
		Handler:      app.routes(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	go gracefulShutdown(srv, app.Logger)

	app.Logger.Infof("Server is running on %s", app.Addr)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		app.Logger.Fatalf("server errror: %v", err)
	}
}

func gracefulShutdown(srv *http.Server, logger logger.Logger) {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-done

	const Duration = 5
	ctx, cancel := context.WithTimeout(context.Background(), Duration*time.Second)
	defer cancel()

	logger.Infof("Shutting down server with %ds timeout", Duration)

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("could not shutdown server: %v", err)
	}
}
