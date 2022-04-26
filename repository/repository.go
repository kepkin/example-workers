package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/golang-migrate/migrate/v4"
	mpgx "github.com/golang-migrate/migrate/v4/database/pgx"

    _ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"

	ipm "github.com/kepkin/interview-price-monitor"
)

const minimalSupportedSchemaVersion = 1
const migrationPath = "file://repository/migrations"
// const migrationPath = "file://migrations"


type PGRepository struct {
	dbpool *sql.DB

	schemaVersion uint
}

func NewRepo(connString string) (ipm.Repository, error) {
	return newRepo(connString)
}

func newRepo(connString string) (*PGRepository, error) {
	db, err := sql.Open("pgx", connString)
	if err != nil {
		return nil, err
	}

	instance, err := mpgx.WithInstance(db, &mpgx.Config{})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(migrationPath, "postgres", instance)
	if err != nil {
		return nil, err		
	}

	schemaVersion, _, err := m.Version()
	if errors.Is(err, migrate.ErrNilVersion) {
		err = m.Up()  // initial migration
		schemaVersion, _, err = m.Version()
		if err != nil {
			return nil, err
		}
	}

	if err != nil {
		return nil, err
	}

	if schemaVersion < minimalSupportedSchemaVersion {
		return nil, fmt.Errorf("miniman supported schema version is: %v", minimalSupportedSchemaVersion)
	}

	return &PGRepository{
		dbpool: db,
		schemaVersion: schemaVersion,
	}, nil
}

func (r *PGRepository) resetDB() error {
	instance, err := mpgx.WithInstance(r.dbpool, &mpgx.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(migrationPath, "postgres", instance)
	if err != nil {
		return err
	}
	err = m.Drop()
	if err != nil {
		return err
	}

	instance, err = mpgx.WithInstance(r.dbpool, &mpgx.Config{})
	if err != nil {
		return err
	}

	m, err = migrate.NewWithDatabaseInstance(migrationPath, "postgres", instance)
	if err != nil {
		return err
	}

	return m.Up()
}

func (r *PGRepository) checkForMinimalSchemaVersion(target uint) error {
	if target > r.schemaVersion {
		return fmt.Errorf("this method is not supported for this schema version")
	}
	return nil
}

func (r *PGRepository) AddMonitor(ctx context.Context, m ipm.MonitorTask) error {
	err := r.checkForMinimalSchemaVersion(1)
	if err != nil {
		return err
	}

	var target_id int64
	err = r.dbpool.QueryRowContext(ctx, "SELECT id FROM targets WHERE name=$1", m.Target).Scan(&target_id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ipm.NewError(ipm.ErrNoSuchTarget, m.Target)
		}
		//TODO:
		return err
	}


	_, err = r.dbpool.ExecContext(ctx, `INSERT INTO monitor_tasks
		(id, start, stop, target_id, frequency, status)
		VALUES($1, $2, $3, $4, $5, $6)`, m.ID, m.Start.UTC(), m.Stop.UTC(), target_id, m.Frequency, m.Status)
	if err != nil {
		var e *pgconn.PgError
		if errors.As(err, &e) {
			if e.Code == "23505" {
				return ipm.NewError(ipm.ErrDuplicateMonitor, m.ID.String())
			}
		} 
	}
	return err
}

//TODO: rename to acquire tasks
func (r *PGRepository) ListMonitorsToStart(ctx context.Context) ([]ipm.MonitorTask, error) {
	err := r.checkForMinimalSchemaVersion(1)
	if err != nil {
		return nil, err
	}

	rows, err := r.dbpool.QueryContext(ctx, `
		WITH cte1 as (UPDATE  monitor_tasks m
			SET status = $2, liveprobe=NOW()
			WHERE start <= NOW() + interval '2 second' 
				AND stop > NOW()
				AND status = $1
			RETURNING
				m.id, m.start, m.stop, m.target_id, m.frequency, m.status)

		SELECT m.id, m.start, m.stop, t.name as target, m.frequency
		FROM cte1 m
		JOIN targets t ON m.target_id = t.id 
		ORDER by start; 
		`, ipm.TaskStatusFree, ipm.TaskStatusRunning)

	if err != nil {
		return nil, err
	}
	defer rows.Close()


	result := make([]ipm.MonitorTask, 0)
	for rows.Next() {
		t := ipm.MonitorTask{}
		
		if err := rows.Scan(&t.ID, &t.Start, &t.Stop, &t.Target, &t.Frequency); err != nil {
			return nil, err
		}
		result = append(result, t)
	}

	return result, err
}

func (r *PGRepository) Put(ctx context.Context, id uuid.UUID, time time.Time, price ipm.Decimal) error {
	err := r.checkForMinimalSchemaVersion(1)
	if err != nil {
		return err
	}

	_, err = r.dbpool.ExecContext(ctx, "INSERT INTO monitor_results (monitor_task_id, time, price) VALUES($1,$2,$3)", id, time, price.String())
	return err
}

func (r *PGRepository) ListMonitorResults(ctx context.Context, id uuid.UUID) (ipm.PriceMap, error) {
	err := r.checkForMinimalSchemaVersion(1)
	if err != nil {
		return nil, err
	}

	rows, err := r.dbpool.QueryContext(ctx, `SELECT m.time, m.price
		FROM monitor_results m 
		WHERE monitor_task_id = $1`,
		id,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()


	result := make(ipm.PriceMap)
	for rows.Next() {
		var t time.Time
		price := ipm.Decimal{}
		
		if err := rows.Scan(&t, &price); err != nil {
			return nil, err
		}
		result[t] = price
	}

	return result, err		
}
