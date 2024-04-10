package internal

import (
	"context"
	"errors"
	"flag"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"os"
)

const (
	runAddressVar        = "RUN_ADDRESS"
	databaseURIVar       = "DATABASE_URI"
	accrualSystemAddrVar = "ACCRUAL_SYSTEM_ADDRESS"
)

type App struct {
	Address           string
	AccrualSysAddress string
	DBPool            *pgxpool.Pool
}

type config struct {
	runAddress        string
	databaseURI       string
	accrualSystemAddr string
}

var Logger zap.SugaredLogger

func InitApp(ctx context.Context) (*App, error) {
	initLogger()
	c := initConfig()
	if err := checkConfig(c); err != nil {
		return nil, err
	}

	dbPool, err := initDB(ctx, c.databaseURI)
	if err != nil {
		return nil, err
	}

	return &App{
		Address:           c.runAddress,
		AccrualSysAddress: c.accrualSystemAddr,
		DBPool:            dbPool,
	}, nil
}

func initConfig() *config {
	c := new(config)

	flag.StringVar(&c.runAddress, "a", "", "server address")
	flag.StringVar(&c.databaseURI, "d", "", "database uri")
	flag.StringVar(&c.accrualSystemAddr, "r", "", "accrual system address")

	flag.Parse()

	if envVar := os.Getenv(runAddressVar); envVar != "" {
		c.runAddress = envVar
	}

	if envVar := os.Getenv(databaseURIVar); envVar != "" {
		c.databaseURI = envVar
	}

	if envVar := os.Getenv(accrualSystemAddrVar); envVar != "" {
		c.accrualSystemAddr = envVar
	}

	return c
}

func initDB(ctx context.Context, DSN string) (*pgxpool.Pool, error) {
	dbConf, err := pgxpool.ParseConfig(DSN)
	if err != nil {
		return nil, err
	}

	dbPool, err := pgxpool.NewWithConfig(ctx, dbConf)
	if err != nil {
		return nil, err
	}

	err = dbPool.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return dbPool, nil
}

func checkConfig(c *config) error {
	if c.runAddress == "" || c.databaseURI == "" || c.accrualSystemAddr == "" {
		return errors.New("please, check configs")
	}

	return nil
}

func initLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	defer func(logger *zap.Logger) {
		err = logger.Sync()
	}(logger)

	if err != nil {
		panic(err)
	}
	Logger = *logger.Sugar()
}
