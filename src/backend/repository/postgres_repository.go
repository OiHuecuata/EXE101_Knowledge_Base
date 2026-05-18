package repository

import (
	"backend/model"
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository interface {
	CreateSession(ctx context.Context, session *model.ChatSession) error
	GetSessions(ctx context.Context) ([]model.ChatSession, error)
	CreateMessage(ctx context.Context, message *model.ChatMessage) error
	GetMessagesBySessionID(ctx context.Context, sessionID uuid.UUID) ([]model.ChatMessage, error)
}

type postgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) PostgresRepository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) CreateSession(ctx context.Context, session *model.ChatSession) error {
	query := `INSERT INTO chat_sessions (id, title, created_at) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(ctx, query, session.ID, session.Title, session.CreatedAt)
	return err
}

func (r *postgresRepository) GetSessions(ctx context.Context) ([]model.ChatSession, error) {
	query := `SELECT id, title, created_at FROM chat_sessions ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []model.ChatSession
	for rows.Next() {
		var s model.ChatSession
		if err := rows.Scan(&s.ID, &s.Title, &s.CreatedAt); err != nil {
			return nil, err
		}
		sessions = append(sessions, s)
	}
	return sessions, nil
}

func (r *postgresRepository) CreateMessage(ctx context.Context, message *model.ChatMessage) error {
	query := `INSERT INTO chat_messages (session_id, role, content, created_at) 
              VALUES ($1, $2, $3, $4) RETURNING id`
	err := r.db.QueryRow(ctx, query, message.SessionID, message.Role, message.Content, message.CreatedAt).Scan(&message.ID)
	return err
}

func (r *postgresRepository) GetMessagesBySessionID(ctx context.Context, sessionID uuid.UUID) ([]model.ChatMessage, error) {
	query := `SELECT id, session_id, role, content, created_at 
              FROM chat_messages WHERE session_id = $1 ORDER BY created_at ASC`
	rows, err := r.db.Query(ctx, query, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []model.ChatMessage
	for rows.Next() {
		var m model.ChatMessage
		if err := rows.Scan(&m.ID, &m.SessionID, &m.Role, &m.Content, &m.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	return messages, nil
}
