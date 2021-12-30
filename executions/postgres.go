// Package executions will manage all executions requirements
package executions

import (
	"context"

	"github.com/Lord-Y/cypress-parallel/commons"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
)

// list will return all executions with range limit settings
func (p *listExecutions) list() (z []dbList, err error) {
	ctx := context.Background()
	db, err := pgxpool.Connect(ctx, commons.BuildDSN())
	if err != nil {
		return
	}
	defer db.Close()

	rows, err := db.Query(
		ctx,
		"SELECT e.*, (SELECT count(e.execution_id) FROM executions e) total, p.project_name FROM executions e LEFT JOIN projects p ON e.project_id = p.project_id ORDER BY e.date DESC OFFSET $1 LIMIT $2",
		p.StartLimit,
		p.EndLimit,
	)
	if err != nil && err.Error() != pgx.ErrNoRows.Error() {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var x dbList
		if err = rows.Scan(
			&x.Execution_id,
			&x.Project_id,
			&x.Branch,
			&x.Execution_status,
			&x.Uniq_id,
			&x.Spec,
			&x.Result,
			&x.Date,
			&x.Execution_error_output,
			&x.Pod_name,
			&x.Pod_cleaned,
			&x.Total,
			&x.Project_name,
		); err != nil {
			return
		}
		z = append(z, x)
	}
	return z, nil
}

// read will return return specific execution content
func (p *readExecutions) read() (z DBRead, err error) {
	ctx := context.Background()
	db, err := pgxpool.Connect(ctx, commons.BuildDSN())
	if err != nil {
		return
	}
	defer db.Close()

	tx, err := db.Begin(ctx)
	if err != nil {
		return
	}
	//golangci-lint fail on this check while the transaction error is checked
	defer tx.Rollback(ctx) //nolint

	err = tx.QueryRow(
		ctx,
		"SELECT e.*, p.project_name FROM executions e LEFT JOIN projects p ON e.project_id = p.project_id WHERE e.execution_id = $1 LIMIT 1",
		p.ExecutionID,
	).Scan(
		&z.Execution_id,
		&z.Project_id,
		&z.Branch,
		&z.Execution_status,
		&z.Uniq_id,
		&z.Spec,
		&z.Result,
		&z.Date,
		&z.Execution_error_output,
		&z.Pod_name,
		&z.Pod_cleaned,
		&z.Project_name,
	)
	if err = tx.Commit(ctx); err != nil {
		return
	}
	return z, nil
}

// updateResult will update execution result in DB
func (p *updateResultExecution) updateResult() (z string, err error) {
	ctx := context.Background()
	db, err := pgxpool.Connect(ctx, commons.BuildDSN())
	if err != nil {
		return
	}
	defer db.Close()

	tx, err := db.Begin(ctx)
	if err != nil {
		return
	}
	//golangci-lint fail on this check while the transaction error is checked
	defer tx.Rollback(ctx) //nolint

	err = tx.QueryRow(
		ctx,
		"UPDATE executions SET result = $1, execution_status = $2, execution_error_output = $3, pod_cleaned = 'true' WHERE uniq_id = $4 AND spec = $5 AND branch = $6 RETURNING pod_name",
		p.Result,
		p.ExecutionStatus,
		p.ExecutionErrorOutput,
		p.UniqID,
		p.Spec,
		p.Branch,
	).Scan(
		&z,
	)
	if err = tx.Commit(ctx); err != nil {
		return
	}
	return z, nil
}

// countExecutions will count number of executions not in specified values
func (p *updateResultExecution) countExecutions() (z dbCountExecutions, err error) {
	ctx := context.Background()
	db, err := pgxpool.Connect(ctx, commons.BuildDSN())
	if err != nil {
		return
	}
	defer db.Close()

	err = db.QueryRow(
		ctx,
		"SELECT pod_name, execution_status FROM executions WHERE uniq_id = $1 AND execution_status = 'RUNNING' AND pod_name = (SELECT pod_name FROM executions WHERE uniq_id = $1 AND spec = $2)",
		p.UniqID,
		p.Spec,
	).Scan(
		&z.Pod_name,
		&z.Execution_status,
	)
	if err != nil && err.Error() != pgx.ErrNoRows.Error() {
		return
	}
	return z, nil
}

// countExecutionsInverted will count number of executions not in specified values
func (p *updateResultExecution) countExecutionsInverted() (z dbCountExecutions, err error) {
	ctx := context.Background()
	db, err := pgxpool.Connect(ctx, commons.BuildDSN())
	if err != nil {
		return
	}
	defer db.Close()

	err = db.QueryRow(
		ctx,
		"SELECT pod_name, execution_status FROM executions WHERE uniq_id = $1 AND spec = $2",
		p.UniqID,
		p.Spec,
	).Scan(
		&z.Pod_name,
		&z.Execution_status,
	)
	if err != nil && err.Error() != pgx.ErrNoRows.Error() {
		return
	}
	return z, nil
}

// GetExecutionIDForUnitTesting will fetch data from db
func GetExecutionIDForUnitTesting() (z DBRead, err error) {
	ctx := context.Background()
	db, err := pgxpool.Connect(ctx, commons.BuildDSN())
	if err != nil {
		return
	}
	defer db.Close()

	tx, err := db.Begin(ctx)
	if err != nil {
		return
	}
	//golangci-lint fail on this check while the transaction error is checked
	defer tx.Rollback(ctx) //nolint

	err = tx.QueryRow(
		ctx,
		"SELECT * FROM executions LIMIT 1",
	).Scan(
		&z.Execution_id,
		&z.Project_id,
		&z.Branch,
		&z.Execution_status,
		&z.Uniq_id,
		&z.Spec,
		&z.Result,
		&z.Date,
		&z.Execution_error_output,
		&z.Pod_name,
		&z.Pod_cleaned,
		&z.Project_name,
	)
	if err = tx.Commit(ctx); err != nil {
		return
	}
	return z, nil
}

// search will return all projects
func (p *searchExecutions) search() (z []dbList, err error) {
	ctx := context.Background()
	db, err := pgxpool.Connect(ctx, commons.BuildDSN())
	if err != nil {
		return
	}
	defer db.Close()

	rows, err := db.Query(
		ctx,
		"SELECT e.*, (SELECT COUNT(e.execution_id) FROM executions e LEFT JOIN projects p ON e.project_id = p.project_id WHERE e.branch ILIKE '%' || $1 || '%' OR e.uniq_id ILIKE '%' || $1 || '%' OR e.spec ILIKE '%' || $1 || '%' OR p.project_name ILIKE '%' || $1 || '%') total, p.project_name FROM executions e LEFT JOIN projects p ON e.project_id = p.project_id WHERE e.branch ILIKE '%' || $1 || '%' OR e.uniq_id ILIKE '%' || $1 || '%' OR e.spec ILIKE '%' || $1 || '%' OR p.project_name ILIKE '%' || $1 || '%' ORDER BY e.date DESC OFFSET $2 LIMIT $3",
		p.Q,
		p.StartLimit,
		p.EndLimit,
	)
	if err != nil && err.Error() != pgx.ErrNoRows.Error() {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var x dbList
		if err = rows.Scan(
			&x.Execution_id,
			&x.Project_id,
			&x.Branch,
			&x.Execution_status,
			&x.Uniq_id,
			&x.Spec,
			&x.Result,
			&x.Date,
			&x.Execution_error_output,
			&x.Pod_name,
			&x.Pod_cleaned,
			&x.Total,
			&x.Project_name,
		); err != nil {
			return
		}
		z = append(z, x)
	}
	return z, nil
}

// uniqId will return all executions of the uniq id provided
func (p *uniqIDExecutions) uniqId() (z []DBRead, err error) {
	ctx := context.Background()
	db, err := pgxpool.Connect(ctx, commons.BuildDSN())
	if err != nil {
		return
	}
	defer db.Close()

	rows, err := db.Query(
		ctx,
		"SELECT e.*, p.project_name FROM executions e LEFT JOIN projects p ON e.project_id = p.project_id WHERE e.uniq_id = $1 ORDER BY e.date DESC",
		p.UniqID,
	)
	if err != nil && err.Error() != pgx.ErrNoRows.Error() {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var x DBRead
		if err = rows.Scan(
			&x.Execution_id,
			&x.Project_id,
			&x.Branch,
			&x.Execution_status,
			&x.Uniq_id,
			&x.Spec,
			&x.Result,
			&x.Date,
			&x.Execution_error_output,
			&x.Pod_name,
			&x.Pod_cleaned,
			&x.Project_name,
		); err != nil {
			return
		}
		z = append(z, x)
	}
	return z, nil
}
