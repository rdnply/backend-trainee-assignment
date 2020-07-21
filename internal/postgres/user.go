package postgres

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/rdnply/backend-trainee-assignment/internal/format"
	"github.com/rdnply/backend-trainee-assignment/internal/user"
)

var _ user.Storage = &UserStorage{}

type UserStorage struct {
	statementStorage

	addStmt  *sql.Stmt
	findStmt *sql.Stmt
}

func NewUserStorage(db *DB) (*UserStorage, error) {
	s := &UserStorage{statementStorage: newStatementsStorage(db)}

	stmts := []stmt{
		{Query: addUserQuery, Dst: &s.addStmt},
		{Query: findUserQuery, Dst: &s.findStmt},
	}

	if err := s.initStatements(stmts); err != nil {
		return nil, errors.Wrap(err, "can't init statements")
	}

	return s, nil
}

const userFields = "username, created_at"
const addUserQuery = "INSERT INTO users(" + userFields + ") VALUES ($1, $2) RETURNING id "

func (s *UserStorage) Add(u *user.User) error {
	if err := s.addStmt.QueryRow(u.Username, format.NewNullTime()).Scan(&u.ID); err != nil {
		return errors.Wrap(err, "can't exec query")
	}

	return nil
}

func scanUser(scanner sqlScanner, u *user.User) error {
	return scanner.Scan(&u.ID, &u.Username, &u.CreatedAt)
}

const findUserQuery = "SELECT id, " + userFields + " FROM users WHERE username=$1"

func (s *UserStorage) Find(username string) (*user.User, error) {
	var u user.User

	row := s.findStmt.QueryRow(username)
	if err := scanUser(row, &u); err != nil {
		if err == sql.ErrNoRows {
			return &u, nil
		}

		return &u, errors.Wrap(err, "can't scan user")
	}

	return &u, nil
}
