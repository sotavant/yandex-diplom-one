package pgsql

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sotavant/yandex-diplom-one/domain"
	"github.com/sotavant/yandex-diplom-one/internal"
	"strings"
)

const withdrawnTableName = "withdrawn"

type WithdrawnRepository struct {
	DBPoll *pgxpool.Pool
}

func NewWithdrawnRepository(ctx context.Context, pool *pgxpool.Pool) (*WithdrawnRepository, error) {
	err := createWithdrawnTable(ctx, pool)
	if err != nil {
		return nil, err
	}

	return &WithdrawnRepository{DBPoll: pool}, nil
}

func (wd *WithdrawnRepository) Store(ctx context.Context, withdrawn domain.Withdrawn) error {
	tx, err := wd.DBPoll.Begin(ctx)
	defer func(ctx context.Context, tx pgx.Tx) {
		err = tx.Rollback(ctx)
		if err != nil {
			internal.Logger.Infow("error in transaction", "err", err)
			panic(err)
		}
	}(ctx, tx)

	query := setWithdrawnTableName(`insert into #T# (order, user_id, sum) values ($1, $2, $3)`)
	userQuery := setUserTableName(`update #T# 
		set withdrawn = withdrawn + $1,
			current = withdrawn - $2
		where id = $3
	`)

	_, err = wd.DBPoll.Query(ctx, query, withdrawn.OrderNum, withdrawn.UserId, withdrawn.Sum)
	if err != nil {
		return err
	}

	_, err = wd.DBPoll.Query(ctx, userQuery, withdrawn.Sum, withdrawn.Sum, withdrawn.UserId)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func createWithdrawnTable(ctx context.Context, pool *pgxpool.Pool) error {
	query := strings.ReplaceAll(`create table if not exists #T#
		(
			id          serial
				constraint withdrawn_pk
					primary key,
			order_num      bigint                  not null
				constraint withdrawn_uq
					unique,
			user_id     bigint                  not null
				constraint withdrawn_users_id_fk
					references public.users,
			sum      float8                 not null,
			processed_at timestamp default now() not null
		);`, "#T#", withdrawnTableName)

	_, err := pool.Exec(ctx, query)

	return err
}

func setWithdrawnTableName(query string) string {
	return strings.ReplaceAll(query, "#T#", withdrawnTableName)
}
