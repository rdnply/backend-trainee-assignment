package postgres

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/rdnply/backend-trainee-assignment/internal/chat"
)

var _ chat.Storage = &ChatStorage{}

type ChatStorage struct {
	statementStorage

	db *DB

	findStmt   *sql.Stmt
	existsStmt *sql.Stmt
}

func NewChatStorage(db *DB) (*ChatStorage, error) {
	s := &ChatStorage{db: db, statementStorage: newStatementsStorage(db)}

	stmts := []stmt{
		{Query: findChatQuery, Dst: &s.findStmt},
		{Query: existsChatQuery, Dst: &s.existsStmt},
	}

	if err := s.initStatements(stmts); err != nil {
		return nil, errors.Wrap(err, "can't init statements")
	}

	return s, nil
}

const addChatQuery = "INSERT INTO chats(name) VALUES ($1) RETURNING chat_id"
const addChatUserRelationQuery = "INSERT INTO chat_user(chat_id, user_id) VALUES($1, $2)"

func (s *ChatStorage) Add(newChat *chat.Chat) error {
	tx, err := s.db.Session.Begin()
	if err != nil {
		return errors.Wrap(err, "can't start transaction")
	}
	defer tx.Rollback() // The rollback will be ignored if the tx has been committed later in the function.

	err = tx.QueryRow(addChatQuery, newChat.Name).Scan(&newChat.ID)
	if err != nil {
		return errors.Wrap(err, "can't add chat")
	}

	stmt, err := tx.Prepare(addChatUserRelationQuery)
	if err != nil {
		return errors.Wrap(err, "can't prepare statement")
	}
	defer stmt.Close()

	for _, userID := range newChat.UsersIDs {
		if _, err := stmt.Exec(newChat.ID, userID); err != nil {
			return errors.Wrap(err, "can't add relation")
		}
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "can't commit transaction")
	}

	return nil
}

const chatFields = "name, created_at"
const findChatQuery = "SELECT chat_id, " + chatFields + " FROM chats WHERE name=$1"

func (s *ChatStorage) Find(name string) (*chat.Chat, error) {
	var c chat.Chat

	row := s.findStmt.QueryRow(name)
	if err := scanChat(row, &c); err != nil {
		if err == sql.ErrNoRows {
			return &c, nil
		}

		return &c, errors.Wrap(err, "can't scan chat")
	}

	return &c, nil
}

func scanChat(scanner sqlScanner, c *chat.Chat) error {
	return scanner.Scan(&c.ID, &c.Name, &c.CreatedAt)
}

const existsChatQuery = "SELECT EXISTS (SELECT chat_id FROM chats WHERE chat_id=$1)"

func (s *ChatStorage) Exists(id int) (bool, error) {
	var exists bool

	if err := s.existsStmt.QueryRow(id).Scan(&exists); err != nil {
		return exists, errors.Wrap(err, "can't exec query")
	}

	return exists, nil
}

func (s *ChatStorage) GetAll(id int) ([]*chat.Chat, error) {
	return nil, nil
}
