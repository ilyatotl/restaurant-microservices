package dish

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

func (r *PostgresRepository) add(ctx context.Context, dish *DishModel) (int64, error) {
	var id int64
	err := r.db.ExecQueryRow(
		ctx,
		`INSERT INTO dishes(name, description, price, quantity) VALUES ($1, $2, $3, $4) RETURNING id`,
		dish.Name,
		dish.Description,
		dish.Price,
		dish.Quantity).Scan(&id)
	return id, err
}

func (r *PostgresRepository) getById(ctx context.Context, id int64) (*DishModel, error) {
	var d DishModel
	err := r.db.Get(ctx, &d, "SELECT * FROM dishes WHERE id = $1", id)
	if err == sql.ErrNoRows {
		return nil, custom_errors.ErrDishNotFound
	}
	return &d, err
}

func (r *PostgresRepository) list(ctx context.Context) ([]*DishModel, error) {
	dishes := make([]*DishModel, 0)
	err := r.db.Select(ctx, &dishes, "SELECT id,name,description,price,quantity,created_at,updated_at FROM dishes")
	return dishes, err
}

func (r *PostgresRepository) update(ctx context.Context, dish *DishModel) (bool, error) {
	result, err := r.db.Exec(ctx,
		"UPDATE dishes SET name = $1, description = $2, price = $3, quantity = $4, updated_at = now() WHERE id = $5",
		dish.Name,
		dish.Description,
		dish.Price,
		dish.Quantity,
		dish.ID)

	return result.RowsAffected() > 0, err
}

func (r *PostgresRepository) delete(ctx context.Context, id int64) (bool, error) {
	result, err := r.db.Exec(ctx,
		"DELETE FROM dishes WHERE id = $1", id)
	return result.RowsAffected() > 0, err
}
