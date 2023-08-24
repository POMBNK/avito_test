package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/POMBNK/avito_test_task/internal/segment"
	"github.com/POMBNK/avito_test_task/pkg/client/postgresql"
	"github.com/POMBNK/avito_test_task/pkg/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"time"
)

type postgresDB struct {
	logs   *logger.Logger
	client postgresql.Client
}

// Create method.
// Create a segment database entries with "name" of segment and "active" boolean flag.
// Active -> true means segment is active otherwise active -> false.
func (d *postgresDB) Create(ctx context.Context, segment segment.Segment) (string, error) {
	var pgErr *pgconn.PgError
	d.logs.Debug("Check if segment already exist")
	err := d.isSegmentExist(ctx, segment)
	if err != nil {
		if errors.As(err, &pgErr) {
			return "", err
		}
		//TODO: Get correct error
		return "", fmt.Errorf("this segment has already exist")
	}

	d.logs.Debug("Creating segment")
	q := `INSERT INTO segment (name, active) VALUES ($1,$2) RETURNING id`
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	err = d.client.QueryRow(ctx, q, segment.Name, segment.Active).Scan(&segment.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", err
		}
		return "", fmt.Errorf("can not create segment due error:%w", err)
	}

	d.logs.Debug("Segment created")
	d.logs.Tracef("id of created segment: %s \n", segment.ID)

	return segment.ID, nil
}

// Delete method.
// Update a segment field "active" to false (0).
// The field is not deleted from table:
//   - Not to corrupt the data in the user entity;
//   - Save statistic data on future.
func (d *postgresDB) Delete(ctx context.Context, segment segment.Segment) error {
	d.logs.Debug("Removing segment...")

	q := `UPDATE segment SET active=0 WHERE name = $1 `
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	res, err := d.client.Exec(ctx, q, segment.Name)
	if err != nil {
		return err
	}

	if res.RowsAffected() != 1 {
		return fmt.Errorf("not found 404") //TODO: apierror.ErrNotFound http code 404
	}

	d.logs.Tracef("Matched and deleted %v segments.\n", res.RowsAffected())
	return nil
}

func (d *postgresDB) isSegmentExist(ctx context.Context, segment segment.Segment) error {
	q := `SELECT id FROM segment WHERE name=$1`
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	err := d.client.QueryRow(ctx, q, segment.Name).Scan(&segment.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return err
	}
	return nil
}

func NewPostgresDB(logs *logger.Logger, client postgresql.Client) segment.Storage {
	return &postgresDB{
		logs:   logs,
		client: client,
	}
}
