package postgresql

import (
	"context"
	"fmt"
	"github.com/POMBNK/avito_test_task/pkg/config"
	"github.com/POMBNK/avito_test_task/pkg/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

func NewClient(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	var pool *pgxpool.Pool
	var err error
	data := cfg.Storage.Postgresql
	dns := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", data.User, data.Password, data.Host, data.Port, data.Database)
	utils.Again(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		pool, err = pgxpool.New(ctx, dns)
		if err != nil {
			return err
		}
		return nil
	}, 5, 5*time.Second)
	return pool, nil
}
