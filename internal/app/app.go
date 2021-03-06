package app

import (
	"io"
	"log"

	"github.com/pkg/errors"
	"github.com/rdnply/backend-trainee-assignment/internal/chat"
	"github.com/rdnply/backend-trainee-assignment/internal/message"
	"github.com/rdnply/backend-trainee-assignment/internal/postgres"
	"github.com/rdnply/backend-trainee-assignment/internal/user"
	"github.com/rdnply/backend-trainee-assignment/pkg/logger"
)

type App struct {
	Addr           string
	UserStorage    user.Storage
	ChatStorage    chat.Storage
	MessageStorage message.Storage
	Logger         logger.Logger
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

	chatStorage, err := postgres.NewChatStorage(db)
	if err != nil {
		return nil, nil, err
	}
	closers["chat_storage"] = chatStorage

	messageStorage, err := postgres.NewMessageStorage(db)
	if err != nil {
		return nil, nil, err
	}
	closers["message_storage"] = messageStorage

	return &App{
		Addr:           addr,
		UserStorage:    userStorage,
		ChatStorage:    chatStorage,
		MessageStorage: messageStorage,
		Logger:         initLogger(),
	}, closers, nil
}

func connectPostgres() (*postgres.DB, error) {
	db, err := postgres.New()
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
