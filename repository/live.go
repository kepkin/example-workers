package repository

import (
	"context"

	"github.com/google/uuid"

	ipm "github.com/kepkin/interview-price-monitor"
)


func (r *PGRepository) CheckLiveProbe(ctx context.Context) (int64, error) {
	err := r.checkForMinimalSchemaVersion(1)
	if err != nil {
		return 0, err
	}

	res, err := r.dbpool.ExecContext(ctx, `
		UPDATE monitor_tasks m
			SET status = $1, liveprobe=NULL
			WHERE liveprobe <= NOW() - interval '10 seconds'
			AND status = $2
		`, ipm.TaskStatusFree, ipm.TaskStatusRunning)

	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}


func (r *PGRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status ipm.TaskStatus) error {
	err := r.checkForMinimalSchemaVersion(1)
	if err != nil {
		return err
	}

	_, err = r.dbpool.ExecContext(ctx, `
		UPDATE monitor_tasks m
			SET status = $1, liveprobe=NOW()
			WHERE id = $2
		`, status, id)

	return err
}
