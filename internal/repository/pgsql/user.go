package pgsql

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sotavant/yandex-diplom-one/domain"
	"strings"
)

const usersTableName = "users"

type UserRepository struct {
	DBPoll *pgxpool.Pool
}

func NewUserRepository(ctx context.Context, pool *pgxpool.Pool) (*UserRepository, error) {
	err := createUsersTable(ctx, pool)
	if err != nil {
		return nil, err
	}

	err = addColumns(ctx, pool)
	if err != nil {
		return nil, err
	}

	return &UserRepository{DBPoll: pool}, nil
}

func (u *UserRepository) GetByID(ctx context.Context, userID int64) (domain.User, error) {
	query := setUserTableName(`select id, login, password, current, withdrawn from #T# where id = $1`)

	return u.getOne(ctx, query, userID)
}

func (u *UserRepository) GetByLogin(ctx context.Context, login string) (domain.User, error) {
	query := setUserTableName(`select id, login, password, current, withdrawn from #T# where login = $1`)

	return u.getOne(ctx, query, login)
}

func (u *UserRepository) Store(ctx context.Context, user domain.User) (int64, error) {
	var id int64
	query := setUserTableName(`insert into #T# (login, password) values ($1, $2) returning id`)

	err := u.DBPoll.QueryRow(ctx, query, user.Login, user.Password).Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil
}

func (u *UserRepository) getOne(ctx context.Context, query string, args ...interface{}) (user domain.User, err error) {
	rows, err := u.DBPoll.Query(ctx, query, args...)
	if err != nil {
		return user, err
	}

	users, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.User])
	if err != nil {
		return user, err
	}

	for _, user = range users {
		return user, nil
	}

	return
}

func createUsersTable(ctx context.Context, pool *pgxpool.Pool) error {
	query := strings.ReplaceAll(`create table if not exists #T#
		(
			id    serial primary key,
			login  varchar not null,
			password varchar not null
		);`, "#T#", usersTableName)

	_, err := pool.Exec(ctx, query)

	return err
}

func addColumns(ctx context.Context, pool *pgxpool.Pool) error {
	query := strings.ReplaceAll(`
		alter table #T# add column if not exists current float8 default 0;
		alter table #T# add column if not exists withdrawn float8 default 0;
	`, "#T#", usersTableName)

	_, err := pool.Exec(ctx, query)

	return err
}

func setUserTableName(query string) string {
	return strings.ReplaceAll(query, "#T#", usersTableName)
}
