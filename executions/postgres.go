// Package executions will manage all executions requirements
package executions

import (
	"database/sql"

	"github.com/Lord-Y/cypress-parallel/commons"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"github.com/syyongx/php2go"
)

// list will return all executions with range limit settings
func (p *listExecutions) list() (z []map[string]interface{}, err error) {
	db, err := sql.Open(
		"postgres",
		commons.BuildDSN(),
	)
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT e.*, (SELECT count(e.execution_id) FROM executions e) total, p.project_name FROM executions e LEFT JOIN projects p ON e.project_id = p.project_id ORDER BY e.date DESC OFFSET $1 LIMIT $2")
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

// read will return return specific execution content
func (p *readExecutions) read() (z []interface{}, err error) {
	db, err := sql.Open(
		"postgres",
		commons.BuildDSN(),
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to DB")
		return z, err
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT e.*, p.project_name FROM executions e LEFT JOIN projects p ON e.project_id = p.project_id WHERE e.execution_id = $1 LIMIT 1")
	if err != nil && err != sql.ErrNoRows {
		return z, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(
		p.ExecutionID,
	)
	if err != nil && err != sql.ErrNoRows {
		return z, err
	}

	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return z, err
	}

	count := len(columnTypes)
	finalRows := []interface{}{}

	for rows.Next() {
		scanArgs := make([]interface{}, count)
		for i, v := range columnTypes {
			switch v.DatabaseTypeName() {
			case "VARCHAR", "TEXT", "UUID", "TIMESTAMP":
				scanArgs[i] = new(sql.NullString)
				break //nolint:gosimple
			case "BOOL":
				scanArgs[i] = new(sql.NullBool)
				break //nolint:gosimple
			case "INT4":
				scanArgs[i] = new(sql.NullInt64)
				break //nolint:gosimple
			default:
				scanArgs[i] = new(sql.NullString)
			}
		}
		err := rows.Scan(scanArgs...)
		if err != nil {
			return z, err
		}

		m := map[string]interface{}{}
		for i, v := range columnTypes {
			if z, ok := (scanArgs[i]).(*sql.NullBool); ok {
				m[v.Name()] = z.Bool
				continue
			}
			if z, ok := (scanArgs[i]).(*sql.NullString); ok {
				m[v.Name()] = z.String
				continue
			}
			if z, ok := (scanArgs[i]).(*sql.NullInt64); ok {
				m[v.Name()] = z.Int64
				continue
			}
			if z, ok := (scanArgs[i]).(*sql.NullFloat64); ok {
				m[v.Name()] = z.Float64
				continue
			}
			if z, ok := (scanArgs[i]).(*sql.NullInt32); ok {
				m[v.Name()] = z.Int32
				continue
			}
			m[v.Name()] = scanArgs[i]
		}
		finalRows = append(finalRows, m)
	}

	if err = rows.Err(); err != nil {
		return z, err
	}
	return finalRows, nil
}

// updateResult will update execution result in DB
func (p *updateResultExecution) updateResult() (z string, err error) {
	db, err := sql.Open(
		"postgres",
		commons.BuildDSN(),
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to DB")
		return z, err
	}
	defer db.Close()

	stmt, err := db.Prepare("UPDATE executions SET result = $1, execution_status = $2, execution_error_output = $3, pod_cleaned = 'true' WHERE uniq_id = $4 AND spec = $5 AND branch = $6 RETURNING pod_name")
	if err != nil && err != sql.ErrNoRows {
		return z, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(
		p.Result,
		php2go.Addslashes(p.ExecutionStatus),
		php2go.Addslashes(p.ExecutionErrorOutput),
		php2go.Addslashes(p.UniqID),
		php2go.Addslashes(p.Spec),
		php2go.Addslashes(p.Branch),
	).Scan(&z)
	if err != nil && err != sql.ErrNoRows {
		return z, err
	}
	return z, nil
}

// countExecutions will count number of executions not in specified values
func (p *updateResultExecution) countExecutions() (z map[string]string, err error) {
	db, err := sql.Open(
		"postgres",
		commons.BuildDSN(),
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to DB")
		return z, err
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT pod_name, execution_status FROM executions WHERE uniq_id = $1 AND execution_status = 'RUNNING' AND pod_name = (SELECT pod_name FROM executions WHERE uniq_id = $1 AND spec = $2)")
	if err != nil && err != sql.ErrNoRows {
		return z, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(
		php2go.Addslashes(p.UniqID),
		php2go.Addslashes(p.Spec),
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

// countExecutionsInverted will count number of executions not in specified values
func (p *updateResultExecution) countExecutionsInverted() (z map[string]string, err error) {
	db, err := sql.Open(
		"postgres",
		commons.BuildDSN(),
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to DB")
		return z, err
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT pod_name, execution_status FROM executions WHERE uniq_id = $1 AND spec = $2")
	if err != nil && err != sql.ErrNoRows {
		return z, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(
		php2go.Addslashes(p.UniqID),
		php2go.Addslashes(p.Spec),
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

// GetExecutionIDForUnitTesting in only for unit testing purpose and will return annotation_id field
func GetExecutionIDForUnitTesting() (z map[string]string, err error) {
	db, err := sql.Open(
		"postgres",
		commons.BuildDSN(),
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to DB")
		return z, err
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT * FROM executions LIMIT 1")
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
func (p *searchExecutions) search() (z []interface{}, err error) {
	db, err := sql.Open(
		"postgres",
		commons.BuildDSN(),
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to DB")
		return z, err
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT e.*, (SELECT count(execution_id) FROM executions WHERE branch LIKE '%' || $1 || '%' OR uniq_id LIKE '%' || $1 || '%' OR spec LIKE '%' || $1 || '%') total, p.project_name FROM executions e LEFT JOIN projects p ON e.project_id = p.project_id WHERE e.branch LIKE '%' || $1 || '%' OR e.uniq_id LIKE '%' || $1 || '%' OR e.spec LIKE '%' || $1 || '%' ORDER BY e.date DESC OFFSET $2 LIMIT $3")
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

	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return z, err
	}

	count := len(columnTypes)
	finalRows := []interface{}{}

	for rows.Next() {
		scanArgs := make([]interface{}, count)
		for i, v := range columnTypes {
			switch v.DatabaseTypeName() {
			case "VARCHAR", "TEXT", "UUID", "TIMESTAMP":
				scanArgs[i] = new(sql.NullString)
				break //nolint:gosimple
			case "BOOL":
				scanArgs[i] = new(sql.NullBool)
				break //nolint:gosimple
			case "INT4":
				scanArgs[i] = new(sql.NullInt64)
				break //nolint:gosimple
			default:
				scanArgs[i] = new(sql.NullString)
			}
		}
		err := rows.Scan(scanArgs...)
		if err != nil {
			return z, err
		}

		m := map[string]interface{}{}
		for i, v := range columnTypes {
			if z, ok := (scanArgs[i]).(*sql.NullBool); ok {
				m[v.Name()] = z.Bool
				continue
			}
			if z, ok := (scanArgs[i]).(*sql.NullString); ok {
				m[v.Name()] = z.String
				continue
			}
			if z, ok := (scanArgs[i]).(*sql.NullInt64); ok {
				m[v.Name()] = z.Int64
				continue
			}
			if z, ok := (scanArgs[i]).(*sql.NullFloat64); ok {
				m[v.Name()] = z.Float64
				continue
			}
			if z, ok := (scanArgs[i]).(*sql.NullInt32); ok {
				m[v.Name()] = z.Int32
				continue
			}
			m[v.Name()] = scanArgs[i]
		}
		finalRows = append(finalRows, m)
	}

	if err = rows.Err(); err != nil {
		return z, err
	}
	return finalRows, nil
}

// uniqId will return all executions of the uniq id provided
func (p *uniqIDExecutions) uniqId() (z []interface{}, err error) {
	db, err := sql.Open(
		"postgres",
		commons.BuildDSN(),
	)
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT e.*, p.project_name FROM executions e LEFT JOIN projects p ON e.project_id = p.project_id WHERE e.uniq_id = $1 ORDER BY e.date DESC")
	if err != nil && err != sql.ErrNoRows {
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(
		php2go.Addslashes(p.UniqID),
	)
	if err != nil && err != sql.ErrNoRows {
		return
	}

	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return z, err
	}

	count := len(columnTypes)
	finalRows := []interface{}{}

	for rows.Next() {
		scanArgs := make([]interface{}, count)
		for i, v := range columnTypes {
			switch v.DatabaseTypeName() {
			case "VARCHAR", "TEXT", "UUID", "TIMESTAMP":
				scanArgs[i] = new(sql.NullString)
				break //nolint:gosimple
			case "BOOL":
				scanArgs[i] = new(sql.NullBool)
				break //nolint:gosimple
			case "INT4":
				scanArgs[i] = new(sql.NullInt64)
				break //nolint:gosimple
			default:
				scanArgs[i] = new(sql.NullString)
			}
		}
		err := rows.Scan(scanArgs...)
		if err != nil {
			return z, err
		}

		m := map[string]interface{}{}
		for i, v := range columnTypes {
			if z, ok := (scanArgs[i]).(*sql.NullBool); ok {
				m[v.Name()] = z.Bool
				continue
			}
			if z, ok := (scanArgs[i]).(*sql.NullString); ok {
				m[v.Name()] = z.String
				continue
			}
			if z, ok := (scanArgs[i]).(*sql.NullInt64); ok {
				m[v.Name()] = z.Int64
				continue
			}
			if z, ok := (scanArgs[i]).(*sql.NullFloat64); ok {
				m[v.Name()] = z.Float64
				continue
			}
			if z, ok := (scanArgs[i]).(*sql.NullInt32); ok {
				m[v.Name()] = z.Int32
				continue
			}
			m[v.Name()] = scanArgs[i]
		}
		finalRows = append(finalRows, m)
	}

	if err = rows.Err(); err != nil {
		return z, err
	}
	return finalRows, nil
}
