package pgsql

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sotavant/yandex-diplom-one/domain"
	"github.com/sotavant/yandex-diplom-one/internal"
	"strings"
)

const OrdersTableName = "orders"

type OrderRepository struct {
	DBPoll *pgxpool.Pool
}

func NewOrderRepository(ctx context.Context, pool *pgxpool.Pool) (*OrderRepository, error) {
	err := createOrdersTable(ctx, pool)

	if err != nil {
		return nil, err
	}

	return &OrderRepository{DBPoll: pool}, nil
}

func (o *OrderRepository) FindByStatus(ctx context.Context, states []string) ([]domain.Order, error) {
	var orders []domain.Order

	query := setOrderTableName(`select id, number, user_id, status, accrual, uploaded_at from #T# where status = any($1)`)

	rows, err := o.DBPoll.Query(ctx, query, states)
	if err != nil {
		return orders, err
	}

	orders, err = pgx.CollectRows(rows, pgx.RowToStructByName[domain.Order])
	if err != nil {
		return make([]domain.Order, 0), err
	}

	return orders, nil
}

func (o *OrderRepository) FindByUser(ctx context.Context, userID int64) ([]domain.Order, error) {
	var orders []domain.Order

	query := setOrderTableName(`select id, number, user_id, status, accrual, uploaded_at from #T# where user_id = $1 order by uploaded_at asc`)

	rows, err := o.DBPoll.Query(ctx, query, userID)
	if err != nil {
		return orders, err
	}

	orders, err = pgx.CollectRows(rows, pgx.RowToStructByName[domain.Order])
	if err != nil {
		return make([]domain.Order, 0), err
	}

	return orders, nil
}

func (o *OrderRepository) GetByNum(ctx context.Context, orderNum string) (domain.Order, error) {
	query := setOrderTableName(`select id, number, user_id, status, accrual, uploaded_at from #T# where number = $1`)

	return o.getOne(ctx, query, orderNum)
}

func (o *OrderRepository) SetAccrual(ctx context.Context, order domain.Order) error {
	tx, err := o.DBPoll.Begin(ctx)
	defer func(tx pgx.Tx, ctx context.Context) {
		err = tx.Rollback(ctx)
		if err != nil {
			internal.Logger.Infow("error in close transaction", "err", err)
		}
	}(tx, ctx)

	query := setOrderTableName(`update #T# 
		set status = $1,
		 	accrual = $2
	where id = $3`)

	userQuery := setUserTableName(`update #T#
		set current = current + $1
	where id = $2`)

	_, err = tx.Exec(ctx, query, order.Status, order.Accrual, order.ID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, userQuery, order.Accrual, order.UserID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (o *OrderRepository) Store(ctx context.Context, order domain.Order) (int64, error) {
	var id int64
	query := setOrderTableName(`insert into #T# (number, user_id, status) values ($1, $2, $3) returning id`)

	err := o.DBPoll.QueryRow(ctx, query, order.Number, order.UserID, order.Status).Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil
}

func (o *OrderRepository) UpdateStatus(ctx context.Context, order domain.Order) error {
	query := setOrderTableName("update #T# set status = $1 where id = $2")

	_, err := o.DBPoll.Exec(ctx, query, order.Status, order.ID)

	if err != nil {
		return err
	}

	return nil
}

func (o *OrderRepository) getOne(ctx context.Context, query string, args ...interface{}) (order domain.Order, err error) {
	rows, err := o.DBPoll.Query(ctx, query, args...)
	if err != nil {
		return order, err
	}

	orders, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.Order])
	if err != nil {
		return order, err
	}

	for _, order = range orders {
		return order, nil
	}

	return
}

func createOrdersTable(ctx context.Context, pool *pgxpool.Pool) error {
	query := strings.ReplaceAll(`create table if not exists #T#
		(
			id          serial
				constraint orders_pk
					primary key,
			number      varchar                  not null
				constraint orders_uq
					unique,
			user_id     bigint                  not null
				constraint orders_users_id_fk
					references public.users,
			status      varchar                 not null,
			accrual     float8,
			uploaded_at timestamp default now() not null
		);`, "#T#", OrdersTableName)

	_, err := pool.Exec(ctx, query)

	if err != nil {
		return err
	}
	return nil
}

func setOrderTableName(query string) string {
	return strings.ReplaceAll(query, "#T#", OrdersTableName)
}
