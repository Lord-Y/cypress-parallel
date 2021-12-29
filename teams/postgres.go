// Package teams will manage all teams requirements
package teams

import (
	"context"
	"fmt"

	"github.com/Lord-Y/cypress-parallel/commons"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
)

// create will insert teams in DB
func (p *teams) create() (z int64, err error) {
	ctx := context.Background()
	db, err := pgxpool.Connect(ctx, commons.BuildDSN())
	if err != nil {
		return
	}
	defer db.Close()

	var count int
	err = db.QueryRow(
		ctx,
		"SELECT COUNT(team_id) FROM teams WHERE team_name = $1",
		p.Name,
	).Scan(
		&count,
	)

	if count > 0 {
		return z, fmt.Errorf("already_exist")
	}

	tx, err := db.Begin(ctx)
	if err != nil {
		return
	}
	//golangci-lint fail on this check while the transaction error is checked
	defer tx.Rollback(ctx) //nolint

	err = tx.QueryRow(
		ctx,
		"INSERT INTO teams(team_name) VALUES($1) RETURNING team_id",
		p.Name,
	).Scan(
		&z,
	)
	if err = tx.Commit(ctx); err != nil {
		return
	}
	return z, nil
}

// read will return all teams with range limit settings
func (p *getTeams) read() (z dbCommon, err error) {
	ctx := context.Background()
	db, err := pgxpool.Connect(ctx, commons.BuildDSN())
	if err != nil {
		return
	}
	defer db.Close()

	err = db.QueryRow(
		ctx,
		"SELECT * FROM teams WHERE team_id = $1 LIMIT 1",
		p.TeamID,
	).Scan(
		&z.Team_id,
		&z.Team_name,
		&z.Date,
	)
	if err != nil && err.Error() != pgx.ErrNoRows.Error() {
		return
	}
	return z, nil
}

// list will return all teams with range limit settings
func (p *listTeams) list() (z []dbList, err error) {
	ctx := context.Background()
	db, err := pgxpool.Connect(ctx, commons.BuildDSN())
	if err != nil {
		return
	}
	defer db.Close()

	rows, err := db.Query(
		ctx,
		"SELECT *, (SELECT COUNT(team_id) FROM teams) total FROM teams ORDER BY date DESC OFFSET $1 LIMIT $2",
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
			&x.Team_id,
			&x.Team_name,
			&x.Date,
			&x.Total,
		); err != nil {
			return
		}
		z = append(z, x)
	}
	return z, nil
}

// all will return all teams
func all() (z []dbList, err error) {
	ctx := context.Background()
	db, err := pgxpool.Connect(ctx, commons.BuildDSN())
	if err != nil {
		return
	}
	defer db.Close()

	rows, err := db.Query(
		ctx,
		"SELECT *, (SELECT count(team_id) FROM teams) total FROM teams ORDER BY date",
	)
	if err != nil && err.Error() != pgx.ErrNoRows.Error() {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var x dbList
		if err = rows.Scan(
			&x.Team_id,
			&x.Team_name,
			&x.Date,
			&x.Total,
		); err != nil {
			return
		}
		z = append(z, x)
	}
	return z, nil
}

// GetTeamIDForUnitTesting in only for unit testing purpose and will return team_id field
func GetTeamIDForUnitTesting() (z dbCommon, err error) {
	ctx := context.Background()
	db, err := pgxpool.Connect(ctx, commons.BuildDSN())
	if err != nil {
		return
	}
	defer db.Close()

	err = db.QueryRow(
		ctx,
		"SELECT * FROM teams LIMIT 1",
	).Scan(
		&z.Team_id,
		&z.Team_name,
		&z.Date,
	)
	if err != nil && err.Error() != pgx.ErrNoRows.Error() {
		return
	}
	return z, nil
}

// update will update team name in DB
func (p *updateTeam) update() (err error) {
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
		"UPDATE teams SET team_name = $1 WHERE team_id = $2",
		p.Name,
		p.TeamID,
	)
	if err = tx.Commit(ctx); err != nil {
		return
	}
	return nil
}

// delete will delete teams in DB
func (p *deleteTeam) delete() (err error) {
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
		"DELETE FROM teams WHERE team_id = $1",
		p.TeamID,
	)
	if err = tx.Commit(ctx); err != nil {
		return
	}
	return nil
}

// search will return all teams
func (p *searchTeams) search() (z []dbList, err error) {
	ctx := context.Background()
	db, err := pgxpool.Connect(ctx, commons.BuildDSN())
	if err != nil {
		return
	}
	defer db.Close()

	rows, err := db.Query(
		ctx,
		"SELECT *, (SELECT count(team_id) FROM teams WHERE team_name ILIKE '%' || $1 || '%') total FROM teams WHERE team_name ILIKE '%' || $1 || '%' ORDER BY date DESC OFFSET $2 LIMIT $3",
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
			&x.Team_id,
			&x.Team_name,
			&x.Date,
			&x.Total,
		); err != nil {
			return
		}
		z = append(z, x)
	}
	return z, nil
}
