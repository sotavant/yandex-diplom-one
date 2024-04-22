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
	/*query1 := `drop table if exists $1`
	query2 := `drop table if exists $1`*/

	_, err := pool.Exec(ctx, query)
	if err != nil {
		return err
	}
	/*
		_, err = pool.Exec(ctx, query1, pgsql.OrdersTableName)
		if err != nil {
			return err
		}

		_, err = pool.Exec(ctx, query2, pgsql.UsersTableName)
		if err != nil {
			return err
		}*/

	return nil
}
