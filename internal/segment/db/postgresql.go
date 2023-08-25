package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/POMBNK/avito_test_task/internal/apierror"
	"github.com/POMBNK/avito_test_task/internal/segment"
	"github.com/POMBNK/avito_test_task/internal/user"
	"github.com/POMBNK/avito_test_task/pkg/client/postgresql"
	"github.com/POMBNK/avito_test_task/pkg/logger"
	"github.com/jackc/pgx/v5"
	"strconv"
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

	d.logs.Debug("Check if segment already exist")
	existID, err := d.isSegmentExist(ctx, segment.Name)
	if err != nil {
		return "", err
	}

	if existID == "" {
		newId, err := d.createSegment(ctx, segment)
		if err != nil {
			return "", err
		}
		return newId, nil
	}

	err = d.makeSegmentActive(ctx, segment)
	if err != nil {
		return "", err
	}

	return existID, nil
}

// Delete method.
// Update a segment field "active" to false (0).
// The field is not deleted from table:
//   - Not to corrupt the data in the user entity;
//   - Save statistic data on future.
func (d *postgresDB) Delete(ctx context.Context, segment segment.Segment) error {
	d.logs.Debug("Removing segment...")

	q := `UPDATE segment SET active='0' WHERE name = $1 `
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	res, err := d.client.Exec(ctx, q, segment.Name)
	if err != nil {
		return err
	}
	if res.RowsAffected() != 1 {
		return apierror.ErrNotFound
	}

	q = `UPDATE user_segment us
			SET active = FALSE, del_at =now()
			FROM segment s
			WHERE us.segment_id = s.id
  			AND us.active = TRUE
  			AND s.name = $1;`
	ctx, cancel = context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	res, err = d.client.Exec(ctx, q, segment.Name)
	if err != nil {
		return err
	}

	return nil
}

// AddUserToSegments Method for adding a user to a segment.
// Accepts a list of (names) of segments to add a user to
func (d *postgresDB) AddUserToSegments(ctx context.Context, segmentsUser segment.SegmentsUsers, segmentName string) error {

	existedSegment, err := d.isSegmentExist(ctx, segmentName)
	if existedSegment == "" {
		return apierror.ErrNotFound
	}

	intUserID, err := strconv.Atoi(segmentsUser.UserID)
	if err != nil {
		return err
	}

	q := fmt.Sprintf(`INSERT INTO user_segment(segment_id, user_id, ACTIVE)
			WITH data AS (SELECT id AS segment_id, %d AS user_id, TRUE AS active FROM segment WHERE name = $1)
			SELECT segment_id, user_id, active FROM data
			WHERE NOT EXISTS (SELECT * FROM user_segment
                  WHERE (user_id = (SELECT user_id from DATA) AND
                         segment_id = (SELECT segment_id from DATA) AND
                         active = TRUE))`, intUserID)
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	_, err = d.client.Exec(ctx, q, segmentName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return apierror.ErrNotFound
		}
		return err
	}

	return nil
}

func (d *postgresDB) IsUserExist(ctx context.Context, segmentsUser segment.SegmentsUsers) error {

	var userUnit user.User
	q := `SELECT name, email FROM users WHERE id=$1`
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	err := d.client.QueryRow(ctx, q, segmentsUser.UserID).Scan(&userUnit.Name, &userUnit.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			d.logs.Debug("User id doesn't exist")
			return apierror.ErrNotFound
		}
		return err
	}
	return nil
}

func (d *postgresDB) isSegmentExist(ctx context.Context, segmentName string) (string, error) {
	var segmentUnit segment.Segment
	q := `SELECT id FROM segment WHERE name=$1`
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	err := d.client.QueryRow(ctx, q, segmentName).Scan(&segmentUnit.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", nil
		}
		return "", err
	}
	return segmentUnit.ID, nil
}

func (d *postgresDB) createSegment(ctx context.Context, segment segment.Segment) (string, error) {
	d.logs.Debug("Creating segment")
	q := `INSERT INTO segment (name, active) VALUES ($1,'1') RETURNING id`
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	err := d.client.QueryRow(ctx, q, segment.Name).Scan(&segment.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", apierror.ErrNotFound
		}
		return "", fmt.Errorf("can not create segment due error:%w", err)
	}

	d.logs.Debug("Segment created")
	d.logs.Tracef("id of created segment: %s \n", segment.ID)

	return segment.ID, nil
}

func (d *postgresDB) makeSegmentActive(ctx context.Context, segment segment.Segment) error {
	d.logs.Debug("Update already existed segment...")
	q := `UPDATE segment SET active='1' WHERE name = $1`

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	res, err := d.client.Exec(ctx, q, segment.Name)
	if err != nil {
		return err
	}

	if res.RowsAffected() != 1 {
		return apierror.ErrNotFound
	}

	d.logs.Tracef("Matched and updated %v segments.\n", res.RowsAffected())
	return nil
}

func NewPostgresDB(logs *logger.Logger, client postgresql.Client) segment.Storage {
	return &postgresDB{
		logs:   logs,
		client: client,
	}
}
