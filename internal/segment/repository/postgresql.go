package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/POMBNK/avito_test_task/internal/apierror"
	"github.com/POMBNK/avito_test_task/internal/segment"
	"github.com/POMBNK/avito_test_task/internal/segment/useCase"
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

func (d *postgresDB) AddToRandomUsers(ctx context.Context, segment segment.Segment, percent int) error {
	//TODO: Add id to user model
	q := fmt.Sprintf(`SELECT id FROM users
			ORDER BY RANDOM()
    		LIMIT (SELECT COUNT(id) FROM users) * 0.%d`, percent)

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	userIDs := make([]int, 0)
	rows, err := d.client.Query(ctx, q)
	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			return err
		}
		userIDs = append(userIDs, id)
	}

	if err = rows.Err(); err != nil {
		return err
	}

	ctx, cancel = context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	for _, userID := range userIDs {

		q = fmt.Sprintf(`INSERT INTO user_segment(segment_id, user_id, active)
			  WITH data AS
			  (
				SELECT id AS segment_id, %d AS user_id, TRUE AS active
	          		FROM segment WHERE name = '%s'
			  )
			  SELECT segment_id, user_id, active
	         FROM data
	         WHERE NOT EXISTS (
	           SELECT * FROM user_segment
	           WHERE user_id = (SELECT user_id FROM data)
	             AND segment_id = (SELECT segment_id FROM data)
	             AND active = TRUE
	         )`, userID, segment.Name)

		_, err = d.client.Exec(ctx, q)
		if err != nil {
			return err
		}
	}

	return nil
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

// TODO: Change docstring
// AddUserToSegments Method for adding a user to a segment.
// Accepts a list of (names) of segments to add a user to
func (d *postgresDB) AddUserToSegments(ctx context.Context, segmentsUser segment.SegmentsUsers, segmentName, deleteAfter string) error {

	existedSegment, err := d.isSegmentExist(ctx, segmentName)
	if existedSegment == "" {
		return apierror.ErrNotFound
	}

	intUserID, err := strconv.Atoi(segmentsUser.UserID)
	if err != nil {
		return err
	}
	// for loving memory...
	//var del_after string
	//var stub string
	//
	//if deleteAfter == "" {
	//	del_after = ""
	//	stub = ""
	//} else {
	//	del_after = ", del_after"
	//	stub = fmt.Sprintf(",CAST('%s' AS TIMESTAMPTZ ) as del_after", deleteAfter)
	//}
	//
	//q := fmt.Sprintf(`
	//		INSERT INTO user_segment(segment_id, user_id, ACTIVE %s)
	//		WITH data AS (SELECT id AS segment_id, %d AS user_id, TRUE AS active
	//			 %s FROM segment WHERE name = '%s')
	//		SELECT segment_id, user_id, active %s
	//		FROM data
	//		WHERE NOT EXISTS (SELECT * FROM user_segment
	//		                           WHERE (user_id = (SELECT user_id from DATA)
	//		                                      AND segment_id = (SELECT segment_id from DATA) AND active = TRUE))`,
	//	del_after, intUserID, stub, segmentName, del_after)

	var q string
	if deleteAfter == "" {
		q = fmt.Sprintf(`INSERT INTO user_segment(segment_id, user_id, active)
			  WITH data AS 
			  (
				SELECT id AS segment_id, %d AS user_id, TRUE AS active    
               		FROM segment WHERE name = '%s'
			  )
			  SELECT segment_id, user_id, active
              FROM data
              WHERE NOT EXISTS (
                SELECT * FROM user_segment
                WHERE user_id = (SELECT user_id FROM data)
                  AND segment_id = (SELECT segment_id FROM data)
                  AND active = TRUE
              )`, intUserID, segmentName)
	} else {
		q = fmt.Sprintf(`INSERT INTO user_segment(segment_id, user_id, active, del_after)
			  WITH data AS 
			  (
				SELECT id AS segment_id, %d AS user_id, TRUE AS active, CAST('%s' AS TIMESTAMP ) AS del_after 
               		FROM segment WHERE name = '%s'
			  )
			  SELECT segment_id, user_id, active, del_after 
              FROM data
              WHERE NOT EXISTS (
                SELECT * FROM user_segment
                WHERE user_id = (SELECT user_id FROM data)
                  AND segment_id = (SELECT segment_id FROM data)
                  AND active = TRUE
              )`, intUserID, deleteAfter, segmentName)
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	_, err = d.client.Exec(ctx, q)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return apierror.ErrNotFound
		}
		return err
	}

	return nil
}

func (d *postgresDB) DeleteSegmentFromUser(ctx context.Context, segmentsUser segment.SegmentsUsers, segmentName string) error {
	q := `UPDATE user_segment us
			SET active = FALSE, del_at =now()
			FROM segment s
			WHERE us.segment_id = s.id
  			AND us.active = TRUE
  			AND us.user_id = $1
  			AND s.name = $2;`

	intUserID, err := strconv.Atoi(segmentsUser.UserID)
	if err != nil {
		return err
	}

	_, err = d.client.Exec(ctx, q, intUserID, segmentName)
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

func (d *postgresDB) GetActiveSegments(ctx context.Context, userID string) ([]segment.ActiveSegments, error) {

	q := `SELECT s.id, s.name
		FROM segment s
		INNER JOIN user_segment us ON us.segment_id = s.id
		WHERE us.user_id = $1 AND us.active = TRUE AND s.active = TRUE`

	intUserID, err := strconv.Atoi(userID)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	rows, err := d.client.Query(ctx, q, intUserID)

	allActSegments := make([]segment.ActiveSegments, 0)
	for rows.Next() {
		var actSegment segment.ActiveSegments
		err = rows.Scan(&actSegment.ID, &actSegment.Name)
		if err != nil {
			return nil, err
		}
		allActSegments = append(allActSegments, actSegment)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return allActSegments, nil
}

func (d *postgresDB) GetUserHistoryOptimized(ctx context.Context, userID, timestampz string) ([]segment.BetterCSVReport, error) {
	// more readable table and optimized query
	q := `SELECT us.user_id, s.name, us.active,us.crt_at,us.del_at
		FROM user_segment us
	 	JOIN segment s ON s.id = us.segment_id
	 	JOIN users u on u.id = us.user_id
		WHERE us.user_id = $1
  		AND  us.crt_at >= $2
  		AND (us.del_at IS NULL OR NOW() >= del_at)
		ORDER BY us.crt_at;`

	intUserID, err := strconv.Atoi(userID)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	rows, err := d.client.Query(ctx, q, intUserID, timestampz)

	reports := make([]segment.BetterCSVReport, 0)
	for rows.Next() {
		var report segment.BetterCSVReport
		err = rows.Scan(&report.UserID, &report.SegmentName, &report.Active, &report.CreatedAt, &report.DeletedAt)
		if err != nil {
			return nil, err
		}
		reports = append(reports, report)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return reports, nil
}

func (d *postgresDB) GetUserHistoryOriginal(ctx context.Context, userID string, timestampz string) ([]segment.CSVReport, error) {
	q := `SELECT us.user_id, s.name,'created' as action,us.crt_at AS date
		FROM user_segment us
    	JOIN segment s ON s.id = us.segment_id
    	JOIN users u ON u.id = us.user_id
		WHERE u.id = $1
  		AND  us.crt_at >= $2
		UNION ALL
		SELECT us.user_id, s.name,'deleted' as action,us.del_at as date
		FROM user_segment us
		JOIN segment s ON s.id = us.segment_id
    	JOIN users u on u.id = us.user_id
		WHERE u.id = $1
  		AND us.active=FALSE
  		AND NOW() >= del_at
		ORDER BY date;`

	intUserID, err := strconv.Atoi(userID)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	rows, err := d.client.Query(ctx, q, intUserID, timestampz)

	reports := make([]segment.CSVReport, 0)
	for rows.Next() {
		var report segment.CSVReport
		err = rows.Scan(&report.UserID, &report.SegmentName, &report.Action, &report.Date)
		if err != nil {
			return nil, err
		}
		reports = append(reports, report)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return reports, nil
}

func (d *postgresDB) CheckSegmentsTTL(ctx context.Context) error {
	q := `UPDATE user_segment
			SET active = FALSE
			WHERE (active = TRUE AND del_after <= now());`

	ctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	_, err := d.client.Exec(ctx, q)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
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

func NewPostgresDB(logs *logger.Logger, client postgresql.Client) useCase.Storage {
	return &postgresDB{
		logs:   logs,
		client: client,
	}
}
