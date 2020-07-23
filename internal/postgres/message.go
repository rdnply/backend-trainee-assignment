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
	getStmt *sql.Stmt
}

func NewMessageStorage(db *DB) (*MessageStorage, error) {
	s := &MessageStorage{statementStorage: newStatementsStorage(db)}

	stmts := []stmt{
		{Query: addMessageQuery, Dst: &s.addStmt},
		{Query: getAllMessagesQuery, Dst: &s.getStmt},
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

const getAllMessagesQuery = "SELECT message_id, chat_id, author_id, text, created_at " +
	"FROM messages WHERE chat_id=$1 ORDER BY created_at DESC;"

func (s *MessageStorage) GetAll(chatID int) ([]*message.Message, error) {
	rows, err := s.getStmt.Query(chatID)
	if err != nil {
		return nil, errors.Wrap(err, "can't get all messages")
	}
	defer rows.Close()

	messages := make([]*message.Message, 0)
	for rows.Next() {
		var m message.Message

		if err := scanMessage(rows, &m); err != nil {
			return nil, errors.Wrap(err, "can't scan message")
		}

		messages = append(messages, &m)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "rows contain error")
	}

	return messages, nil
}

func scanMessage(scanner sqlScanner, m *message.Message) error {
	return scanner.Scan(&m.ID, &m.ChatID, &m.AuthorID, &m.Text, &m.CreatedAt)
}
