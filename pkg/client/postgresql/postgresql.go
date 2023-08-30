package postgresql

import (
	"context"
	"fmt"
	"github.com/POMBNK/avito_test_task/pkg/config"
	"github.com/POMBNK/avito_test_task/pkg/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"time"
)

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

func NewClient(ctx context.Context, maxAttempts int, cfg *config.Config) (*pgxpool.Pool, error) {
	var pool *pgxpool.Pool
	var err error
	data := cfg.Storage.Postgresql
	dns := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", data.User, data.Password, data.Host, data.Port, data.Database)
	err = utils.Again(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		pool, err = pgxpool.New(ctx, dns)
		if err != nil {
			return err
		}

		err = pool.Ping(ctx)
		if err != nil {
			return err
		}
		return nil
	}, maxAttempts, 5*time.Second)
	if err != nil {
		log.Fatal("tries limit exceeded")
	}
	return pool, nil
}
