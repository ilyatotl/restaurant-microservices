package dish

import (
	"database/sql"
	"time"
)

type DishModel struct {
	ID          int64        `db:"id" json:"id"`
	Name        string       `db:"name" json:"name"`
	Description string       `db:"description" json:"description"`
	Price       int          `db:"price" json:"price"`
	Quantity    int          `db:"quantity" json:"quantity"`
	CreatedAt   time.Time    `db:"created_at" json:"-"`
	UpdatedAt   sql.NullTime `db:"updated_at" json:"-"`
}

type DishRequest struct {
	ID       int64 `json:"id"`
	Quantity int   `json:"quantity"`
}
