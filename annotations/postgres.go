// Package annotations will manage all annotations requirements
package annotations

import (
	"database/sql"

	"github.com/Lord-Y/cypress-parallel-api/commons"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"github.com/syyongx/php2go"
)

// create will insert annotations in DB
func (p *annotation) create() (z int64, err error) {
	db, err := sql.Open(
		"postgres",
		commons.BuildDSN(),
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to DB")
		return z, err
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO annotations(key, value, project_id) VALUES($1, $2, $3) RETURNING annotation_id")
	if err != nil && err != sql.ErrNoRows {
		return z, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(
		php2go.Addslashes(p.Key),
		php2go.Addslashes(p.Value),
		p.ProjectID,
	).Scan(&z)
	if err != nil && err != sql.ErrNoRows {
		return z, err
	}
	return z, nil
}

// update will update annotations in DB
func (p *updateAnnotation) update() (err error) {
	db, err := sql.Open(
		"postgres",
		commons.BuildDSN(),
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to DB")
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare("UPDATE annotations SET key = $1, value = $2, project_id = $3 WHERE annotation_id = $4")
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	defer stmt.Close()
	err = stmt.QueryRow(
		php2go.Addslashes(p.Key),
		php2go.Addslashes(p.Value),
		p.ProjectID,
		p.AnnotationID,
	).Scan()
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	return nil
}

// list will return all annotations with range limit settings
func (p *listAnnotations) list() (z []map[string]interface{}, err error) {
	db, err := sql.Open(
		"postgres",
		commons.BuildDSN(),
	)
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT a.*, (SELECT count(annotation_id) FROM annotations) total, p.project_name FROM annotations a LEFT JOIN projects p ON a.project_id = p.project_id ORDER BY a.date DESC OFFSET $1 LIMIT $2")
	if err != nil && err != sql.ErrNoRows {
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(
		p.StartLimit,
		p.EndLimit,
	)
	if err != nil && err != sql.ErrNoRows {
		return
	}

	columns, err := rows.Columns()
	if err != nil {
		return
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
			return
		}
		var value string
		sub := make(map[string]interface{})
		for i, col := range values {
			if col == nil {
				value = ""
			} else {
				value = php2go.Stripslashes(string(col))
			}
			sub[columns[i]] = value
		}
		m = append(m, sub)
	}
	if err = rows.Err(); err != nil {
		return
	}
	return m, nil
}

// read will return a single annotation with specified id
func (p *getAnnotations) read() (z map[string]string, err error) {
	db, err := sql.Open(
		"postgres",
		commons.BuildDSN(),
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to DB")
		return z, err
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT a.*, p.project_name FROM annotations a LEFT JOIN projects p ON a.project_id = p.project_id WHERE a.annotation_id = $1 LIMIT 1")
	if err != nil && err != sql.ErrNoRows {
		return z, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(
		p.AnnotationID,
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
				value = ""
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

// delete will delete annotations in DB
func (p *deleteAnnotation) delete() (err error) {
	db, err := sql.Open(
		"postgres",
		commons.BuildDSN(),
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to DB")
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare("DELETE FROM annotations WHERE annotation_id = $1")
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	defer stmt.Close()
	err = stmt.QueryRow(
		p.AnnotationID,
	).Scan()
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	return nil
}

// GetAnnotationIDForUnitTesting in only for unit testing purpose and will return annotation_id field
func GetAnnotationIDForUnitTesting() (z map[string]string, err error) {
	db, err := sql.Open(
		"postgres",
		commons.BuildDSN(),
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to DB")
		return z, err
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT * FROM annotations LIMIT 1")
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
				value = ""
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

// search will return all projects
func (p *searchAnnotations) search() (z []map[string]interface{}, err error) {
	db, err := sql.Open(
		"postgres",
		commons.BuildDSN(),
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to DB")
		return z, err
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT a.*, (SELECT count(annotation_id) FROM annotations WHERE key LIKE '%' || $1 || '%' OR value LIKE '%' || $1 || '%') total, p.project_name FROM annotations a LEFT JOIN projects p ON a.project_id = p.project_id WHERE a.key LIKE '%' || $1 || '%' OR a.value LIKE '%' || $1 || '%' ORDER BY a.date DESC OFFSET $2 LIMIT $3")
	if err != nil && err != sql.ErrNoRows {
		return z, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(
		p.Q,
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
				value = ""
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
