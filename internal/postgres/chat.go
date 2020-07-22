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

	findStmt *sql.Stmt
}

func NewChatStorage(db *DB) (*ChatStorage, error) {
	s := &ChatStorage{db: db, statementStorage: newStatementsStorage(db)}

	stmts := []stmt{
		{Query: findChatQuery, Dst: &s.findStmt},
	}

	if err := s.initStatements(stmts); err != nil {
		return nil, errors.Wrap(err, "can't init statements")
	}

	return s, nil
}

const addChatQuery = "INSERT INTO chats(name) VALUES ($1) RETURNING chat_id"
const addChatUserRelationQuery = "INSERT INTO chat_user(chat_id, user_id) VALUES($1, $2)"

func (s *ChatStorage) Add(chatName string, userIDs []int) (int, error) {
	tx, err := s.db.Session.Begin()
	if err != nil {
		return 0, errors.Wrap(err, "can't start transaction")
	}
	defer tx.Rollback() // The rollback will be ignored if the tx has been committed later in the function.

	var chatID int
	err = tx.QueryRow(addChatQuery, chatName).Scan(&chatID)
	if err != nil {
		return 0, errors.Wrap(err, "can't add chat")
	}

	stmt, err := tx.Prepare(addChatUserRelationQuery)
	if err != nil {
		return 0, errors.Wrap(err, "can't prepare statement")
	}
	defer stmt.Close()

	for _, userID := range userIDs {
		if _, err := stmt.Exec(chatID, userID); err != nil {
			return 0, errors.Wrap(err, "can't add relation")
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, errors.Wrap(err, "can't commit transaction")
	}

	return chatID, nil
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
