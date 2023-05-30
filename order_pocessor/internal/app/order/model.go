package order

import (
	"database/sql"
	"time"
)

type OrderModel struct {
	ID        int64        `db:"id"`
	UserID    int64        `db:"user_id"`
	Status    string       `db:"status"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
}
