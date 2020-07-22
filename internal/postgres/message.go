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

func (s *MessageStorage) Add(m *message.Message) error {
	if err := s.addStmt.QueryRow(m.ChatID, m.AuthorID, m.Text).Scan(&m.ID); err != nil {
		return errors.Wrap(err, "can't exec query")
	}

	return nil
}

func (m *MessageStorage) GetAll() ([]*message.Message, error) {
	panic("not implemented") // TODO: Implement
}
