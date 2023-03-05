package common

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type DependencyInjectionContainer struct {
	Pool   *pgxpool.Pool
	Config *Config
}

func NewDIContainer() (*DependencyInjectionContainer, error) {
	config := ConfigFromEnv()

	pool, err := connectToDatabase(config.Database.URI())
	if err != nil {
		return nil, err
	}

	return &DependencyInjectionContainer{
		Pool:   pool,
		Config: config,
	}, nil
}

func connectToDatabase(uri string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), uri)
	if err != nil {
		return nil, errors.New("cannot initialise postgres connection pool")
	}
	if err := pool.Ping(context.Background()); err != nil {
		return nil, errors.New("unable to connect to database")
	}

	return pool, nil
}
