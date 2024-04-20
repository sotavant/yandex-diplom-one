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

func (o *WithdrawnRepository) FindByUser(ctx context.Context, userId int64) ([]domain.Withdrawn, error) {
	var wds []domain.Withdrawn

	query := setWithdrawnTableName(`select * from #T# where user_id = $1 order by processed_at asc`)

	rows, err := o.DBPoll.Query(ctx, query, userId)
	if err != nil {
		return wds, err
	}

	wds, err = pgx.CollectRows(rows, pgx.RowToStructByName[domain.Withdrawn])
	if err != nil {
		return make([]domain.Withdrawn, 0), err
	}

	return wds, nil
}

func (wd *WithdrawnRepository) FindOne(ctx context.Context, orderNum string) (domain.Withdrawn, error) {
	query := setWithdrawnTableName(`select * from #T# where order_num = $1`)

	return wd.getOne(ctx, query, orderNum)
}

func (wd *WithdrawnRepository) Store(ctx context.Context, withdrawn domain.Withdrawn) error {
	tx, err := wd.DBPoll.Begin(ctx)
	defer func(tx pgx.Tx, ctx context.Context) {
		err = tx.Rollback(ctx)
		if err != nil {
			internal.Logger.Infow("error in close transaction", "err", err)
		}
	}(tx, ctx)

	query := setWithdrawnTableName(`insert into #T# (order_num, user_id, sum) values ($1, $2, $3)`)
	userQuery := setUserTableName(`update #T# 
		set withdrawn = withdrawn + $1,
			current = current - $2
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
			order_num      varchar                  not null
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

func (o *WithdrawnRepository) getOne(ctx context.Context, query string, args ...interface{}) (wd domain.Withdrawn, err error) {
	rows, err := o.DBPoll.Query(ctx, query, args...)
	if err != nil {
		return
	}

	wds, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.Withdrawn])
	if err != nil {
		return
	}

	for _, wd = range wds {
		return
	}

	return
}
