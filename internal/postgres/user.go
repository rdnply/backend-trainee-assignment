package postgres

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/rdnply/backend-trainee-assignment/internal/user"
)

var _ user.Storage = &UserStorage{}

type UserStorage struct {
	statementStorage

	addStmt    *sql.Stmt
	findStmt   *sql.Stmt
	existsStmt *sql.Stmt
}

func NewUserStorage(db *DB) (*UserStorage, error) {
	s := &UserStorage{statementStorage: newStatementsStorage(db)}

	stmts := []stmt{
		{Query: addUserQuery, Dst: &s.addStmt},
		{Query: findUserQuery, Dst: &s.findStmt},
		{Query: existsUserQuery, Dst: &s.existsStmt},
	}

	if err := s.initStatements(stmts); err != nil {
		return nil, errors.Wrap(err, "can't init statements")
	}

	return s, nil
}

const addUserQuery = "INSERT INTO users(username) VALUES ($1) RETURNING user_id "

func (s *UserStorage) Add(u *user.User) error {
	if err := s.addStmt.QueryRow(u.Username).Scan(&u.ID); err != nil {
		return errors.Wrap(err, "can't exec query")
	}

	return nil
}

const userFields = "username, created_at"
const findUserQuery = "SELECT user_id, " + userFields + " FROM users WHERE username=$1"

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

func scanUser(scanner sqlScanner, u *user.User) error {
	return scanner.Scan(&u.ID, &u.Username, &u.CreatedAt)
}

const existsUserQuery = "SELECT EXISTS (SELECT user_id FROM users WHERE user_id=$1)"

func (s *UserStorage) Exists(id int) (bool, error) {
	var exists bool

	if err := s.existsStmt.QueryRow(id).Scan(&exists); err != nil {
		return exists, errors.Wrap(err, "can't exec query")
	}

	return exists, nil
}

func (s *UserStorage) AllExists(ids []int) (bool, int, error) {
	for _, id := range ids {
		exists, err := s.Exists(id)
		if err != nil {
			return false, 0, err
		}
		if !exists {
			return false, id, nil
		}
	}

	return true, 0, nil
}
