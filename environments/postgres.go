// Package environments will manage all environments requirements
package environments

import (
	"context"

	"github.com/Lord-Y/cypress-parallel/commons"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
)

// create will insert environments in DB
func (p *environment) create() (z int64, err error) {
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
		"INSERT INTO environments(key, value, project_id) VALUES($1, $2, $3) RETURNING environment_id",
		p.Key,
		p.Value,
		p.ProjectID,
	).Scan(
		&z,
	)
	if err = tx.Commit(ctx); err != nil {
		return
	}
	return z, nil
}

// update will update environments in DB
func (p *updateEnvironment) update() (err error) {
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
		"UPDATE environments SET key = $1, value = $2, project_id = $3 WHERE environment_id = $4",
		p.Key,
		p.Value,
		p.ProjectID,
		p.EnvironmentID,
	)
	if err = tx.Commit(ctx); err != nil {
		return
	}
	return nil
}

// list will return all environments with range limit settings
func (p *listEnvironments) list() (z []dbList, err error) {
	ctx := context.Background()
	db, err := pgxpool.Connect(ctx, commons.BuildDSN())
	if err != nil {
		return
	}
	defer db.Close()

	rows, err := db.Query(
		ctx,
		"SELECT e.*, (SELECT COUNT(environment_id) FROM environments) total, p.project_name FROM environments e LEFT JOIN projects p ON e.project_id = p.project_id ORDER BY e.date DESC OFFSET $1 LIMIT $2",
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
			&x.Environment_id,
			&x.Key,
			&x.Value,
			&x.Project_id,
			&x.Date,
			&x.Total,
			&x.Project_name,
		); err != nil {
			return
		}
		z = append(z, x)
	}
	return z, nil
}

// read will return all environments with range limit settings
func (p *getEnvironments) read() (z dbRead, err error) {
	ctx := context.Background()
	db, err := pgxpool.Connect(ctx, commons.BuildDSN())
	if err != nil {
		return
	}
	defer db.Close()

	err = db.QueryRow(
		ctx,
		"SELECT e.*, p.project_name FROM environments e LEFT JOIN projects p ON e.project_id = p.project_id WHERE e.environment_id = $1 LIMIT 1",
		p.EnvironmentID,
	).Scan(
		&z.Environment_id,
		&z.Key,
		&z.Value,
		&z.Project_id,
		&z.Date,
		&z.Project_name,
	)
	if err != nil && err.Error() != pgx.ErrNoRows.Error() {
		return
	}
	return z, nil
}

// delete will delete environments in DB
func (p *deleteEnvironment) delete() (err error) {
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
		"DELETE FROM environments WHERE environment_id = $1",
		p.EnvironmentID,
	)
	if err = tx.Commit(ctx); err != nil {
		return
	}
	return nil
}

// GetEnvironmentIDForUnitTesting in only for unit testing purpose and will return environment_id field
func GetEnvironmentIDForUnitTesting() (z dbCommon, err error) {
	ctx := context.Background()
	db, err := pgxpool.Connect(ctx, commons.BuildDSN())
	if err != nil {
		return
	}
	defer db.Close()

	err = db.QueryRow(
		ctx,
		"SELECT * FROM environments LIMIT 1",
	).Scan(
		&z.Environment_id,
		&z.Key,
		&z.Value,
		&z.Project_id,
		&z.Date,
	)
	if err != nil && err.Error() != pgx.ErrNoRows.Error() {
		return
	}
	return z, nil
}

// search will return all projects
func (p *searchEnvironments) search() (z []dbList, err error) {
	ctx := context.Background()
	db, err := pgxpool.Connect(ctx, commons.BuildDSN())
	if err != nil {
		return
	}
	defer db.Close()

	rows, err := db.Query(
		ctx,
		"SELECT e.*, (SELECT COUNT(e.environment_id) FROM environments e LEFT JOIN projects p ON e.project_id = p.project_id WHERE e.key LIKE '%' || $1 || '%' OR e.value ILIKE '%' || $1 || '%' OR p.project_name ILIKE '%' || $1 || '%') total, p.project_name FROM environments e LEFT JOIN projects p ON e.project_id = p.project_id WHERE e.key ILIKE '%' || $1 || '%' OR e.value ILIKE '%' || $1 || '%' OR p.project_name ILIKE '%' || $1 || '%' ORDER BY e.date DESC OFFSET $2 LIMIT $3",
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
			&x.Environment_id,
			&x.Key,
			&x.Value,
			&x.Project_id,
			&x.Date,
			&x.Total,
			&x.Project_name,
		); err != nil {
			return
		}
		z = append(z, x)
	}
	return z, nil
}

// listByProjectID handle requirements to list environments by project id
func (p *listEnvironmentsByProjectID) listByProjectID() (z []dbRead, err error) {
	ctx := context.Background()
	db, err := pgxpool.Connect(ctx, commons.BuildDSN())
	if err != nil {
		return
	}
	defer db.Close()

	rows, err := db.Query(
		ctx,
		"SELECT e.*, p.project_name FROM environments e LEFT JOIN projects p ON e.project_id = p.project_id WHERE e.project_id = $1 ORDER BY e.date DESC",
		p.ProjectID,
	)
	if err != nil && err.Error() != pgx.ErrNoRows.Error() {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var x dbRead
		if err = rows.Scan(
			&x.Environment_id,
			&x.Key,
			&x.Value,
			&x.Project_id,
			&x.Date,
			&x.Project_name,
		); err != nil {
			return
		}
		z = append(z, x)
	}
	return z, nil
}
