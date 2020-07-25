package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq" // init postgres driver
	"github.com/pkg/errors"
)

type DB struct {
	Session *sql.DB
}

func New() (*DB, error) {
	url, err := getDBConfig()
	if err != nil {
		return nil, errors.Wrap(err, "can't get env vars with db configuration")
	}

	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, errors.Wrap(err, "can't open connection to postgres")
	}

	return &DB{
		Session: db,
	}, nil
}

func getDBConfig() (string, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	url := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	return url, nil
}

func (d *DB) CheckConnection() error {
	var err error

	const maxAttempts = 3
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if err = d.Session.Ping(); err == nil {
			break
		}

		nextAttemptWait := time.Duration(attempt) * time.Second
		log.Printf("attempt %d: can't establish a connection with the db, wait for %v: %s",
			attempt,
			nextAttemptWait,
			err,
		)
		time.Sleep(nextAttemptWait)
	}

	return errors.Wrap(err, "can't connect to db")
}

func (d *DB) Close() error {
	if err := d.Session.Close(); err != nil {
		return errors.Wrap(err, "can't close db")
	}

	return nil
}

type sqlScanner interface {
	Scan(dest ...interface{}) error
}
