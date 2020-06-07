package psql

import (
	"context"
	"database/sql"
	"time"

	"github.com/andreyAKor/otus_hw/hw12_13_14_15_calendar/internal/repository"

	// Import Postgres sql driver
	_ "github.com/jackc/pgx/v4/stdlib"
)

var _ repository.DBEventsRepo = (*Repo)(nil)

type Repo struct {
	db *sql.DB
}

// Create connection pool.
func (r *Repo) Connect(ctx context.Context, dsn string) (err error) {
	r.db, err = sql.Open("pgx", dsn)
	if err != nil {
		return
	}
	return r.db.PingContext(ctx)
}

// Close connection pool.
func (r *Repo) Close() error {
	return r.db.Close()
}

// Add new event.
func (r *Repo) Create(ctx context.Context, ev repository.Event) (eventID int64, err error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return eventID, err
	}
	defer func() {
		err = tx.Rollback()
	}()

	if err = r.searchDublicate(
		ctx, tx,
		`SELECT id FROM events WHERE "date" = $1 AND user_id = $2`,
		ev.Date.Format("2006-01-02 15:04:00 -0700"),
		ev.UserID,
	); err != nil {
		return eventID, err
	}

	query := `INSERT INTO
events (title, "date", duration, descr, user_id, duration_start)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id`
	err = tx.QueryRowContext(
		ctx, query,
		ev.Title,
		ev.Date.Format("2006-01-02 15:04:00 -0700"),
		ev.Duration,
		ev.Descr,
		ev.UserID,
		ev.DurationStart,
	).Scan(&eventID)
	if err != nil {
		return eventID, err
	}

	if err = tx.Commit(); err != nil {
		return eventID, err
	}

	return eventID, nil
}

// Update event by id.
func (r *Repo) Update(ctx context.Context, id int64, ev repository.Event) (err error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return
	}
	defer func() {
		err = tx.Rollback()
	}()

	if err = r.searchDublicate(
		ctx, tx,
		`SELECT id FROM events WHERE "date" = $1 AND user_id = $2 AND id != $3`,
		ev.Date.Format("2006-01-02 15:04:00 -0700"),
		ev.UserID,
		id,
	); err != nil {
		return
	}

	query := `UPDATE events
SET title = $1, "date" = $2, duration = $3, descr = $4, duration_start = $5, updated_at = $6
WHERE id = $7`
	res, err := tx.ExecContext(
		ctx, query,
		ev.Title,
		ev.Date,
		ev.Duration,
		ev.Descr,
		ev.DurationStart,
		"now()",
		id,
	)

	ra, err := res.RowsAffected()
	if err != nil {
		return
	}
	if ra == 0 {
		return repository.ErrNotFound
	}

	if err = tx.Commit(); err != nil {
		return
	}

	return
}

// Delete event by id.
func (r *Repo) Delete(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM events WHERE id = $1`, id)
	if err != nil {
		return err
	}

	ra, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if ra == 0 {
		return repository.ErrNotFound
	}

	return nil
}

// Delete last year's events.
func (r *Repo) DeleteOld(ctx context.Context) (err error) {
	yearAgo := time.Now().AddDate(-1, 0, 0)
	_, err = r.db.ExecContext(ctx, `DELETE FROM events WHERE "date" <= $1`, yearAgo)
	return
}

// Get list events on date.
func (r *Repo) GetListByDate(ctx context.Context, date time.Time) ([]repository.Event, error) {
	rows, err := r.db.QueryContext(
		ctx, `SELECT * FROM events WHERE "date" >= $1 AND "date" < $2`,
		date.Format("2006-01-02"),
		date.AddDate(0, 0, 1).Format("2006-01-02"),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.prepareRowsEvents(rows)
}

// Get list events on week by start date.
func (r *Repo) GetListByWeek(ctx context.Context, start time.Time) ([]repository.Event, error) {
	rows, err := r.db.QueryContext(
		ctx, `SELECT * FROM events WHERE "date" >= $1 AND "date" < $2`,
		start.Format("2006-01-02"),
		start.AddDate(0, 0, 8).Format("2006-01-02"),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.prepareRowsEvents(rows)
}

// Get list events on month by start date.
func (r *Repo) GetListByMonth(ctx context.Context, start time.Time) (events []repository.Event, err error) {
	rows, err := r.db.QueryContext(
		ctx, `SELECT * FROM events WHERE "date" >= $1 AND "date" < $2`,
		start.Format("2006-01-02"),
		start.AddDate(0, 0, 31).Format("2006-01-02"),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.prepareRowsEvents(rows)
}

func (r *Repo) prepareRowsEvents(rows *sql.Rows) ([]repository.Event, error) {
	var events []repository.Event

	for rows.Next() {
		var (
			ev                      repository.Event
			updatedAt               sql.NullTime
			duration, durationStart sql.NullInt64
			descr                   sql.NullString
		)

		if err := rows.Scan(
			&ev.ID,
			&ev.Title,
			&ev.Date,
			&duration,
			&descr,
			&ev.UserID,
			&durationStart,
			&ev.CreatedAt,
			&updatedAt,
		); err != nil {
			return nil, err
		}

		if duration.Valid {
			ev.Duration = time.Duration(duration.Int64)
		}
		if descr.Valid {
			ev.Descr = &descr.String
		}
		if durationStart.Valid {
			ev.DurationStart = new(time.Duration)
			*ev.DurationStart = time.Duration(durationStart.Int64)
		}
		if updatedAt.Valid {
			ev.UpdatedAt = new(time.Time)
			*ev.UpdatedAt = updatedAt.Time
		}
		events = append(events, ev)
	}

	return events, rows.Err()
}

// Search dublicate events by date.
func (r *Repo) searchDublicate(ctx context.Context, tx *sql.Tx, query string, args ...interface{}) error {
	var id int64

	err := tx.QueryRowContext(ctx, query, args...).Scan(&id)
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		return err
	}

	return repository.ErrTimeBusy
}
