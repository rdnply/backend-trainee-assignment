package app

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/pkg/errors"
	"github.com/rdnply/backend-trainee-assignment/internal/chat"
	"github.com/rdnply/backend-trainee-assignment/internal/message"
	"github.com/rdnply/backend-trainee-assignment/internal/postgres"
	"github.com/rdnply/backend-trainee-assignment/internal/user"
	"github.com/rdnply/backend-trainee-assignment/pkg/logger"
)

type App struct {
	addr           string
	userStorage    user.Storage
	chatStorage    chat.Storage
	messageStorage message.Storage
	logger         logger.Logger
}

func New(addr string) (*App, map[string]io.Closer, error) {
	closers := make(map[string]io.Closer)

	db, err := connectPostgres()
	if err != nil {
		return nil, nil, err
	}
	closers["postgres"] = db

	userStorage, err := postgres.NewUserStorage(db)
	if err != nil {
		return nil, nil, err
	}
	closers["user_storage"] = userStorage

	return &App{
		addr:        addr,
		userStorage: userStorage,
		logger:      initLogger(),
	}, closers, nil
}

func connectPostgres() (*postgres.DB, error) {
	_ = os.Chdir("../..")

	pwd, err := os.Getwd()
	if err != nil {
		return nil, errors.Wrap(err, "can't get path")
	}

	db, err := postgres.New(fmt.Sprintf("%s/config/postgres.json", pwd))
	if err != nil {
		return nil, errors.Wrap(err, "can't create database instance")
	}

	err = db.CheckConnection()
	if err != nil {
		return nil, errors.Wrap(err, "can't connect to database")
	}

	return db, nil
}

func initLogger() logger.Logger {
	config := logger.Configuration{
		EnableConsole:     true,
		ConsoleLevel:      logger.Debug,
		ConsoleJSONFormat: true,
		EnableFile:        true,
		FileLevel:         logger.Info,
		FileJSONFormat:    true,
		FileLocation:      "log.log",
	}

	logger, err := logger.New(config, logger.InstanceZapLogger)
	if err != nil {
		log.Fatal("could not instantiate logger: ", err)
	}

	return logger
}
