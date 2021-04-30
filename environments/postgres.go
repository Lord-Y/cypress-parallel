// Package environments will manage all environments requirements
package environments

import (
	"database/sql"

	"github.com/Lord-Y/cypress-parallel-api/commons"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"github.com/syyongx/php2go"
)

// selectBeforeAct will check if an insert or update must be done
func (p *Environments) selectBeforeAct(id int) (z map[string]string, err error) {
	db, err := sql.Open(
		"postgres",
		commons.BuildDSN(),
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to DB")
		return z, err
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT count(environment_id) total FROM environments WHERE key = $1 AND value = $2 AND project_id = $3 LIMIT 1")
	if err != nil && err != sql.ErrNoRows {
		return z, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(
		php2go.Addslashes(p.environment[id].Key),
		php2go.Addslashes(p.environment[id].Value),
		p.ProjectID,
	)
	if err != nil && err != sql.ErrNoRows {
		return z, err
	}

	columns, err := rows.Columns()
	if err != nil {
		return z, err
	}

	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	m := make(map[string]string)
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return
		}
		var value string
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = php2go.Stripslashes(string(col))
			}
			m[columns[i]] = value
		}
	}
	if err = rows.Err(); err != nil {
		return z, err
	}
	return m, nil
}

// create will insert environments in DB
func (p *Environments) create(id int) (err error) {
	db, err := sql.Open(
		"postgres",
		commons.BuildDSN(),
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to DB")
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO environments(key, value, project_id) VALUES($1, $2, $3)")
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	err = stmt.QueryRow(
		php2go.Addslashes(p.environment[id].Key),
		php2go.Addslashes(p.environment[id].Value),
		p.ProjectID,
	).Scan()
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	defer stmt.Close()
	return nil
}

// update will insert environments in DB
func (p *Environments) update(id int) (err error) {
	db, err := sql.Open(
		"postgres",
		commons.BuildDSN(),
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to DB")
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare("UPDATE environments SET key = $1, value = $2, project_id = $3 WHERE environment_id = $4")
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	err = stmt.QueryRow(
		php2go.Addslashes(p.environment[id].Key),
		php2go.Addslashes(p.environment[id].Value),
		p.ProjectID,
		p.environment[id].Environment_ID,
	).Scan()
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	defer stmt.Close()
	return nil
}

// read will return all environments with range limit settings
func (p *GetEnvironments) read() (z []map[string]interface{}, err error) {
	db, err := sql.Open(
		"postgres",
		commons.BuildDSN(),
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to DB")
		return z, err
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT *, (SELECT count(environment_id) FROM environments) total FROM environments OFFSET $1 LIMIT $2")
	if err != nil && err != sql.ErrNoRows {
		return z, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(
		p.StartLimit,
		p.EndLimit,
	)
	if err != nil && err != sql.ErrNoRows {
		return z, err
	}

	columns, err := rows.Columns()
	if err != nil {
		return z, err
	}

	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	m := make([]map[string]interface{}, 0)
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return z, err
		}
		var value string
		sub := make(map[string]interface{})
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = php2go.Stripslashes(string(col))
			}
			sub[columns[i]] = value
		}
		m = append(m, sub)
	}
	if err = rows.Err(); err != nil {
		return z, err
	}
	return m, nil
}

// delete will delete environments in DB
func (p *DeleteEnvironment) delete() (err error) {
	db, err := sql.Open(
		"postgres",
		commons.BuildDSN(),
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to DB")
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare("DELETE FROM environments WHERE environment_id = $1")
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	err = stmt.QueryRow(
		p.EnvironmentID,
	).Scan()
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	defer stmt.Close()
	return nil
}

// GetEnvironmentIDForUnitTesting in only for unit testing purpose and will return environment_id field
func GetEnvironmentIDForUnitTesting() (z map[string]string, err error) {
	db, err := sql.Open(
		"postgres",
		commons.BuildDSN(),
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to DB")
		return z, err
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT environment_id FROM environments LIMIT 1")
	if err != nil && err != sql.ErrNoRows {
		return z, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil && err != sql.ErrNoRows {
		return z, err
	}

	columns, err := rows.Columns()
	if err != nil {
		return z, err
	}

	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	m := make(map[string]string)
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return
		}
		var value string
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = php2go.Stripslashes(string(col))
			}
			m[columns[i]] = value
		}
	}
	if err = rows.Err(); err != nil {
		return z, err
	}
	return m, nil
}
