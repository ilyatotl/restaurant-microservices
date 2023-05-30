package sessions

import (
	"authorization/internal/app/custom_errors"
	"context"
	"database/sql"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type dbOps interface {
	Select(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error)
	ExecQueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row
	GetPool(ctx context.Context) *pgxpool.Pool
}

type PostgresRepository struct {
	db dbOps
}

func NewRepository(db dbOps) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) add(ctx context.Context, s *SessionModel) (int64, error) {
	var id int64
	err := r.db.ExecQueryRow(
		ctx,
		`INSERT INTO sessions(user_id, session_token, expires_at) VALUES ($1, $2, $3) RETURNING id`,
		s.UserID,
		s.SessionToken,
		s.ExpiresAt).Scan(&id)
	return id, err
}

func (r *PostgresRepository) get(ctx context.Context, session string) (*SessionModel, error) {
	var s SessionModel
	err := r.db.Get(ctx, &s, "SELECT * FROM sessions WHERE session_token = $1", session)
	if err == sql.ErrNoRows {
		return nil, custom_errors.ErrSessionNotFound
	}
	return &s, err
}

func (r *PostgresRepository) update(ctx context.Context, s *SessionModel) {
	r.db.Exec(ctx,
		"UPDATE sessions SET session_token = $1 WHERE id = $2",
		s.SessionToken,
		s.ID)
}
