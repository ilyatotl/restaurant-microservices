package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	host     = "172.17.0.1"
	port     = 5433
	user     = "ilya"
	password = "test"
	dbname   = "processor"
)

func generateDsn() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
}

func NewDB(ctx context.Context) (*Database, error) {
	dsn := generateDsn()
	pool, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		return nil, err
	}
	return NewDatabase(pool), nil
}
