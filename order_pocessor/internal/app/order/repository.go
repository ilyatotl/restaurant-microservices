package order

import (
	"context"
	"database/sql"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"order_pocessor/internal/app/custom_errors"
)

const Pending = "Pending"

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

func (r *PostgresRepository) add(ctx context.Context, o *OrderModel) (int64, error) {
	var id int64
	err := r.db.ExecQueryRow(
		ctx,
		`INSERT INTO orders(user_id, status) VALUES ($1, $2) RETURNING id`,
		o.UserID,
		o.Status).Scan(&id)
	return id, err
}

func (r *PostgresRepository) getByID(ctx context.Context, id int64) (*OrderModel, error) {
	var o OrderModel
	err := r.db.Get(ctx, &o, "SELECT * FROM orders WHERE id = $1", id)
	if err == sql.ErrNoRows {
		return nil, custom_errors.ErrOrderNotFound
	}
	return &o, err
}

func (r *PostgresRepository) getPendingOrder(ctx context.Context) (*OrderModel, error) {
	var o OrderModel
	err := r.db.Get(ctx, &o, "SELECT * FROM orders WHERE status = $1 LIMIT 1", Pending)
	if err == sql.ErrNoRows {
		return nil, custom_errors.ErrOrderNotFound
	}
	return &o, err
}

func (r *PostgresRepository) changeStatus(ctx context.Context, id int64, status string) (bool, error) {
	res, err := r.db.Exec(ctx, "UPDATE orders SET status = $1 WHERE id = $2", status, id)
	return res.RowsAffected() > 0, err
}
