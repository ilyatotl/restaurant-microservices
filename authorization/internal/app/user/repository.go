package user

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

func (r *PostgresRepository) add(ctx context.Context, u *UserModel) (int64, error) {
	var id int64
	err := r.db.ExecQueryRow(
		ctx,
		`INSERT INTO users(username, email, password_hash, role) VALUES ($1, $2, $3, $4) RETURNING id`,
		u.Username,
		u.Email,
		u.PasswordHash,
		u.Role).Scan(&id)
	return id, err
}

func (r *PostgresRepository) get(ctx context.Context, email string) (*UserModel, error) {
	var u UserModel
	err := r.db.Get(ctx, &u, "SELECT * FROM users WHERE email = $1", email)
	if err == sql.ErrNoRows {
		return nil, custom_errors.ErrUserNotFound
	}
	return &u, err
}

func (r *PostgresRepository) getByID(ctx context.Context, id int64) (*UserDTO, error) {
	var u UserDTO
	err := r.db.Get(ctx, &u, "SELECT id,username,email,role FROM users WHERE id = $1", id)
	if err == sql.ErrNoRows {
		return nil, custom_errors.ErrUserNotFound
	}
	return &u, err
}
