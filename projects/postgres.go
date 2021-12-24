// Package projects will manage all projects requirements
package projects

import (
	"context"

	"github.com/Lord-Y/cypress-parallel/commons"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
)

// create will insert projects in DB
func (p *projects) create() (z int64, err error) {
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
		"INSERT INTO projects(project_name, team_id, repository, branch, specs, scheduling, scheduling_enabled, max_pods, cypress_docker_version, username, password, browser, config_file, timeout) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) RETURNING project_id",
		p.Name,
		p.TeamID,
		p.Repository,
		p.Branch,
		p.Specs,
		p.Scheduling,
		p.SchedulingEnabled,
		p.MaxPods,
		p.CypressDockerVersion,
		p.Username,
		p.Password,
		p.Browser,
		p.ConfigFile,
		p.Timeout,
	).Scan(
		&z,
	)
	if err = tx.Commit(ctx); err != nil {
		return
	}
	return z, nil
}

// read will return a single project with specified id
func (p *getProjects) read() (z dbRead, err error) {
	ctx := context.Background()
	db, err := pgxpool.Connect(ctx, commons.BuildDSN())
	if err != nil {
		return
	}
	defer db.Close()

	err = db.QueryRow(
		ctx,
		"SELECT p.*, t.team_name FROM projects p LEFT JOIN teams t ON p.team_id = t.team_id WHERE project_id = $1 LIMIT 1",
		p.ProjectID,
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
		&z.Team_name,
	)
	if err != nil && err.Error() != pgx.ErrNoRows.Error() {
		return
	}
	return z, nil
}

// list will return all projects with range limit settings
func (p *listProjects) list() (z []dbList, err error) {
	ctx := context.Background()
	db, err := pgxpool.Connect(ctx, commons.BuildDSN())
	if err != nil {
		return
	}
	defer db.Close()

	rows, err := db.Query(
		ctx,
		"SELECT p.*, (SELECT COUNT(project_id) FROM projects) total, t.team_name FROM projects p LEFT JOIN teams t ON p.team_id = t.team_id ORDER BY p.date DESC OFFSET $1 LIMIT $2",
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
			&x.Project_id,
			&x.Project_name,
			&x.Date,
			&x.Team_id,
			&x.Repository,
			&x.Branch,
			&x.Specs,
			&x.Scheduling,
			&x.Scheduling_enabled,
			&x.Max_pods,
			&x.Cypress_docker_version,
			&x.Timeout,
			&x.Username,
			&x.Password,
			&x.Browser,
			&x.Config_file,
			&x.Total,
			&x.Team_name,
		); err != nil {
			return
		}
		z = append(z, x)
	}
	return z, nil
}

// all will return all projects
func all() (z []dbAll, err error) {
	ctx := context.Background()
	db, err := pgxpool.Connect(ctx, commons.BuildDSN())
	if err != nil {
		return
	}
	defer db.Close()

	rows, err := db.Query(
		ctx,
		"SELECT project_id, project_name, branch, cypress_docker_version, browser, config_file, (SELECT count(project_id) FROM projects) total FROM projects ORDER BY date",
	)
	if err != nil && err.Error() != pgx.ErrNoRows.Error() {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var x dbAll
		if err = rows.Scan(
			&x.Project_id,
			&x.Project_name,
			&x.Branch,
			&x.Cypress_docker_version,
			&x.Browser,
			&x.Config_file,
			&x.Total,
		); err != nil {
			return
		}
		z = append(z, x)
	}
	return z, nil
}

// GetProjectIDForUnitTesting in only for unit testing purpose and will return project_id and team_id field
func GetProjectIDForUnitTesting() (z dbCommon, err error) {
	ctx := context.Background()
	db, err := pgxpool.Connect(ctx, commons.BuildDSN())
	if err != nil {
		return
	}
	defer db.Close()

	err = db.QueryRow(
		ctx,
		"SELECT * FROM projects LIMIT 1",
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

// update will update environments in DB
func (p *updateProjects) update() (err error) {
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
		"UPDATE projects SET project_name = $1, team_id = $2, repository = $3, branch = $4, specs = $5, scheduling = $6, scheduling_enabled = $7, max_pods = $8, cypress_docker_version = $9, username = $10, password = $11, browser = $12, config_file = $13, timeout = $14 WHERE project_id = $15",
		p.Name,
		p.TeamID,
		p.Repository,
		p.Branch,
		p.Specs,
		p.Scheduling,
		p.SchedulingEnabled,
		p.MaxPods,
		p.CypressDockerVersion,
		p.Username,
		p.Password,
		p.Browser,
		p.ConfigFile,
		p.Timeout,
		p.ProjectID,
	)
	if err = tx.Commit(ctx); err != nil {
		return
	}
	return nil
}

// delete will delete projects in DB
func (p *deleteProject) delete() (err error) {
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
		"DELETE FROM projects WHERE project_id = $1",
		p.ProjectID,
	)
	if err = tx.Commit(ctx); err != nil {
		return
	}
	return nil
}

// search will return all projects
func (p *searchProjects) search() (z []dbList, err error) {
	ctx := context.Background()
	db, err := pgxpool.Connect(ctx, commons.BuildDSN())
	if err != nil {
		return
	}
	defer db.Close()

	rows, err := db.Query(
		ctx,
		"SELECT p.*, (SELECT count(project_id) FROM projects WHERE project_name ILIKE '%' || $1 || '%') total, t.team_name FROM projects p LEFT JOIN teams t ON p.team_id = t.team_id WHERE p.project_name ILIKE '%' || $1 || '%' ORDER BY p.date DESC OFFSET $2 LIMIT $3",
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
			&x.Project_id,
			&x.Project_name,
			&x.Date,
			&x.Team_id,
			&x.Repository,
			&x.Branch,
			&x.Specs,
			&x.Scheduling,
			&x.Scheduling_enabled,
			&x.Max_pods,
			&x.Cypress_docker_version,
			&x.Timeout,
			&x.Username,
			&x.Password,
			&x.Browser,
			&x.Config_file,
			&x.Total,
			&x.Team_name,
		); err != nil {
			return
		}
		z = append(z, x)
	}
	return z, nil
}
