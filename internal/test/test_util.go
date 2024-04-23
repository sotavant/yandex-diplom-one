package test

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
)

func InitConnection(ctx context.Context) (*pgxpool.Pool, error) {
	dns := os.Getenv("DATABASE_DSN_TEST")
	if dns == "" {
		return nil, nil
	}

	dbConn, err := pgxpool.New(ctx, dns)

	if err != nil {
		return nil, err
	}

	err = dbConn.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return dbConn, nil
}

func CleanData(ctx context.Context, pool *pgxpool.Pool) error {
	query := `drop table if exists withdrawn, orders, users`

	_, err := pool.Exec(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func AddTestUser(ctx context.Context, login, pass string, pool *pgxpool.Pool) error {
	query := `insert into #T# (login, password) values ($1, $2)`

	_, err := pool.Exec(ctx, query, login, pass)

	if err != nil {
		return err
	}

	return nil
}
