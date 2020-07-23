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

const getAllChatsQuery = "SELECT c.chat_id, c.name, c.created_at FROM chats c " +
	"JOIN chat_user cu ON cu.chat_id = c.chat_id " +
	"JOIN " +
	"(SELECT chat_id, MAX(created_at) as max_time " +
	"FROM messages GROUP BY chat_id " +
	")AS m ON m.chat_id = c.chat_id " +
	"WHERE cu.user_id = $1 " +
	"ORDER BY max_time;"

const getAllUserIDsQuery = "SELECT u.user_id FROM users u " +
	"JOIN chat_user cu ON cu.user_id = u.user_id " +
	"WHERE cu.chat_id = $1;"

func (s *ChatStorage) GetAll(userID int) ([]*chat.Chat, error) {
	tx, err := s.db.Session.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "can't start transaction")
	}
	defer tx.Rollback() // The rollback will be ignored if the tx has been committed later in the function.

	rows, err := tx.Query(getAllChatsQuery, userID)
	if err != nil {
		return nil, errors.Wrap(err, "can't get all chats")
	}
	defer rows.Close()

	chats := make([]*chat.Chat, 0)
	for rows.Next() {
		var c chat.Chat
		if err := scanChat(rows, &c); err != nil {
			return nil, errors.Wrap(err, "can't scan row with chat")
		}

		chats = append(chats, &c)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "rows contain error")
	}

	stmt, err := tx.Prepare(getAllUserIDsQuery)
	if err != nil {
		return nil, errors.Wrap(err, "can't prepare statement")
	}
	defer stmt.Close()

	chats, err = addUserIDs(stmt, chats)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, errors.Wrap(err, "can't commit transaction")
	}

	return chats, nil

}

func addUserIDs(stmt *sql.Stmt, chats []*chat.Chat) ([]*chat.Chat, error) {
	for _, chat := range chats {
		rows, err := stmt.Query(chat.ID)
		if err != nil {
			return nil, errors.Wrap(err, "can't get all users in chat")
		}

		ids := make([]int, 0)
		for rows.Next() {
			var id int
			if err := rows.Scan(&id); err != nil {
				return nil, errors.Wrap(err, "can't scan row with id")
			}

			ids = append(ids, id)
		}

		if err = rows.Err(); err != nil {
			return nil, errors.Wrap(err, "rows contain error")
		}

		chat.UsersIDs = ids
	}

	return chats, nil
}
