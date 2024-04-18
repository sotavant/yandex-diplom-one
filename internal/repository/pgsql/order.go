package pgsql

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sotavant/yandex-diplom-one/domain"
	"strings"
)

const ordersTableName = "orders"

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

func (o *OrderRepository) FindByUser(ctx context.Context, userId int64) ([]domain.Order, error) {
	var orders []domain.Order

	query := setOrderTableName(`select * from #T# where user_id = $1 order by uploaded_at asc`)

	rows, err := o.DBPoll.Query(ctx, query, userId)
	if err != nil {
		return orders, err
	}

	orders, err = pgx.CollectRows(rows, pgx.RowToStructByName[domain.Order])
	if err != nil {
		return make([]domain.Order, 0), err
	}

	return orders, nil
}

func (o *OrderRepository) GetByNum(ctx context.Context, orderNum int64) (domain.Order, error) {
	query := setOrderTableName(`select * from #T# where number = $1`)

	return o.getOne(ctx, query, orderNum)
}

func (o *OrderRepository) Store(ctx context.Context, order domain.Order) (int64, error) {
	var id int64
	query := setOrderTableName(`insert into #T# (number, user_id, status) values ($1, $2, $3) returning id`)

	err := o.DBPoll.QueryRow(ctx, query, order.Number, order.UserId, order.Status).Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil
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
			number      bigint                  not null
				constraint orders_uq
					unique,
			user_id     bigint                  not null
				constraint orders_users_id_fk
					references public.users,
			status      varchar                 not null,
			accrual     bigint,
			uploaded_at timestamp default now() not null
		);`, "#T#", ordersTableName)

	_, err := pool.Exec(ctx, query)

	return err
}

func setOrderTableName(query string) string {
	return strings.ReplaceAll(query, "#T#", ordersTableName)
}
