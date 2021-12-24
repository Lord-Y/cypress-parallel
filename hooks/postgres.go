// Package hooks will manage all hooks requirements
package hooks

import (
	"context"

	"github.com/Lord-Y/cypress-parallel/commons"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
)

// getProjectInfos collect requirements to start the unit testing
func (p *plain) getProjectInfos() (z dbProject, err error) {
	ctx := context.Background()
	db, err := pgxpool.Connect(ctx, commons.BuildDSN())
	if err != nil {
		return
	}
	defer db.Close()

	err = db.QueryRow(
		ctx,
		"SELECT * FROM projects WHERE project_name = $1 LIMIT 1",
		p.ProjectName,
	).Scan(
		&z.Project_id,
		&z.Project_name,
		&z.Date,
		&z.Team_id,
		&z.Repository,
		&z.Branch,
		&z.Specs,
		&z.Scheduling,
		&z.Scheduling_enabled,
		&z.Max_pods,
		&z.Cypress_docker_version,
		&z.Timeout,
		&z.Username,
		&z.Password,
		&z.Browser,
		&z.Config_file,
	)
	if err != nil && err.Error() != pgx.ErrNoRows.Error() {
		return
	}
	return z, nil
}

// create will insert executions in DB
func (p *execution) create() (z int64, err error) {
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
		"INSERT INTO executions(project_id, branch, execution_status, uniq_id, spec, result) VALUES($1, $2, $3, $4, $5, $6) RETURNING execution_id",
		p.projectID,
		p.branch,
		p.executionStatus,
		p.uniqID,
		p.spec,
		p.result,
	).Scan(
		&z,
	)
	if err = tx.Commit(ctx); err != nil {
		return
	}
	return z, nil
}

// getProjectAnnotations collect annotatons
func (p *dbProject) getProjectAnnotations() (z []dbProjectAnnotation, err error) {
	ctx := context.Background()
	db, err := pgxpool.Connect(ctx, commons.BuildDSN())
	if err != nil {
		return
	}
	defer db.Close()

	rows, err := db.Query(
		ctx,
		"SELECT * FROM annotations WHERE project_id = $1 LIMIT 1",
		p.Project_id,
	)
	if err != nil && err.Error() != pgx.ErrNoRows.Error() {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var x dbProjectAnnotation
		if err = rows.Scan(
			&x.Annotation_id,
			&x.Key,
			&x.Value,
			&x.Project_id,
			&x.Date,
		); err != nil {
			return
		}
		z = append(z, x)
	}
	return z, nil
}

// getProjectEnvironments collect requirements
func (p *dbProject) getProjectEnvironments() (z []dbProjectEnvironment, err error) {
	ctx := context.Background()
	db, err := pgxpool.Connect(ctx, commons.BuildDSN())
	if err != nil {
		return
	}
	defer db.Close()

	rows, err := db.Query(
		ctx,
		"SELECT * FROM environments WHERE project_id = $1",
		p.Project_id,
	)
	if err != nil && err.Error() != pgx.ErrNoRows.Error() {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var x dbProjectEnvironment
		if err = rows.Scan(
			&x.Environment_id,
			&x.Key,
			&x.Value,
			&x.Project_id,
			&x.Date,
		); err != nil {
			return
		}
		z = append(z, x)
	}
	return z, nil
}

// update will update pod_name field in DB
func (p *updatePodName) update() (err error) {
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

	_, err = tx.Exec(
		ctx,
		"UPDATE executions SET pod_name = $1, execution_status = 'RUNNING' WHERE uniq_id = $2 AND spec = $3",
		p.podName,
		p.uniqID,
		p.spec,
	)
	if err = tx.Commit(ctx); err != nil {
		return
	}
	return nil
}

// executionStatus get executions by status
func executionStatus(execution_status string) (z []dbProjectUniqIDExecution, err error) {
	ctx := context.Background()
	db, err := pgxpool.Connect(ctx, commons.BuildDSN())
	if err != nil {
		return
	}
	defer db.Close()

	rows, err := db.Query(
		ctx,
		"SELECT DISTINCT uniq_id FROM executions WHERE execution_status = $1",
		execution_status,
	)
	if err != nil && err.Error() != pgx.ErrNoRows.Error() {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var x dbProjectUniqIDExecution
		if err = rows.Scan(
			&x.Uniq_id,
		); err != nil {
			return
		}
		z = append(z, x)
	}
	return z, nil
}

// pgqueued get queued executions
func pgqueued(uniqId string) (z []dbPGQueue, err error) {
	ctx := context.Background()
	db, err := pgxpool.Connect(ctx, commons.BuildDSN())
	if err != nil {
		return
	}
	defer db.Close()

	rows, err := db.Query(
		ctx,
		"SELECT e.*, p.project_name FROM executions e LEFT JOIN projects p ON e.project_id = p.project_id WHERE e.execution_status = 'QUEUED' AND e.uniq_id = $1",
		uniqId,
	)
	if err != nil && err.Error() != pgx.ErrNoRows.Error() {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var x dbPGQueue
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

// countExecutions will count number of executions not in specified values
func countExecutions(uniq_id string) (z dbCountExecutions, err error) {
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
		"SELECT COUNT(execution_id) FROM executions WHERE uniq_id = $1 AND execution_status = 'RUNNING'",
		uniq_id,
	).Scan(
		&z.Count,
	)
	if err = tx.Commit(ctx); err != nil {
		return
	}
	return z, nil
}
