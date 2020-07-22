package postgres

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/rdnply/backend-trainee-assignment/internal/message"
)

var _ message.Storage = &MessageStorage{}

type MessageStorage struct {
	statementStorage

	addStmt *sql.Stmt
}

func NewMessageStorage(db *DB) (*MessageStorage, error) {
	s := &MessageStorage{statementStorage: newStatementsStorage(db)}

	stmts := []stmt{
		{Query: addMessageQuery, Dst: &s.addStmt},
	}

	if err := s.initStatements(stmts); err != nil {
		return nil, errors.Wrap(err, "can't init statements")
	}

	return s, nil
}

const addMessageQuery = "INSERT INTO messages(chat_id, author_id, text) VALUES ($1,$2,$3) RETURNING message_id "

func (s *MessageStorage) Add(chatID, authorID int, text string) (int, error) {
	var messageID int

	if err := s.addStmt.QueryRow(chatID, authorID, text).Scan(&messageID); err != nil {
		return 0, errors.Wrap(err, "can't exec query")
	}

	return messageID, nil
}

func (m *MessageStorage) GetAll() ([]*message.Message, error) {
	panic("not implemented") // TODO: Implement
}
