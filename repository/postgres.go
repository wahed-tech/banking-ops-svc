package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wahed-tech/banking-ops-svc/conf"
	"github.com/wahed-tech/go-lib/v2/errors"
)

func NewPostgresPool(ctx context.Context, c conf.PostgresConf) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, c.URL())
	if err != nil {
		return nil, errors.Wrap(err).Here()
	}

	return pool, nil
}
