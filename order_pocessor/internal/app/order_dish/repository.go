package order_dish

import (
	"context"
	"database/sql"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"order_pocessor/internal/app/custom_errors"
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

func (r *PostgresRepository) add(ctx context.Context, o *OrderDishModel) (int64, error) {
	var id int64
	err := r.db.ExecQueryRow(
		ctx,
		`INSERT INTO order_dish(order_id, dish_id, quantity, price) VALUES ($1, $2, $3, $4) RETURNING id`,
		o.OrderID,
		o.DishID,
		o.Quantity,
		o.Price).Scan(&id)
	return id, err
}

func (r *PostgresRepository) getByID(ctx context.Context, id int64) (*OrderDishModel, error) {
	var o OrderDishModel
	err := r.db.Get(ctx, &o, "SELECT * FROM order_dish WHERE id = $1", id)
	if err == sql.ErrNoRows {
		return nil, custom_errors.ErrOrderDishNotFound
	}
	return &o, err
}
