// Package annotations will manage all annotations requirements
package annotations

import (
	"context"

	"github.com/Lord-Y/cypress-parallel/commons"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/lib/pq"
)

// create will insert annotations in DB
func (p *annotation) create() (z int64, err error) {
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
		"INSERT INTO annotations(key, value, project_id) VALUES($1,$2,$3) RETURNING annotation_id",
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

// update will update annotations in DB
func (p *updateAnnotation) update() (err error) {
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
		"UPDATE annotations SET key = $1, value = $2, project_id = $3 WHERE annotation_id = $4",
		p.Key,
		p.Value,
		p.ProjectID,
		p.AnnotationID,
	)
	if err = tx.Commit(ctx); err != nil {
		return
	}
	return nil
}

// list will return all annotations with range limit settings
func (p *listAnnotations) list() (z []dbList, err error) {
	ctx := context.Background()
	db, err := pgxpool.Connect(ctx, commons.BuildDSN())
	if err != nil {
		return
	}
	defer db.Close()

	rows, err := db.Query(
		ctx,
		"SELECT a.*, (SELECT COUNT(annotation_id) FROM annotations) total, p.project_name FROM annotations a LEFT JOIN projects p ON a.project_id = p.project_id ORDER BY a.date DESC OFFSET $1 LIMIT $2",
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
			&x.Annotation_id,
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

// read will return a single annotation with specified id
func (p *getAnnotations) read() (z dbRead, err error) {
	ctx := context.Background()
	db, err := pgxpool.Connect(ctx, commons.BuildDSN())
	if err != nil {
		return
	}
	defer db.Close()

	err = db.QueryRow(
		ctx,
		"SELECT a.*, p.project_name FROM annotations a LEFT JOIN projects p ON a.project_id = p.project_id WHERE a.annotation_id = $1 LIMIT 1",
		p.AnnotationID,
	).Scan(
		&z.Annotation_id,
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

// delete will delete annotations in DB
func (p *deleteAnnotation) delete() (err error) {
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
		"DELETE FROM annotations WHERE annotation_id = $1",
		p.AnnotationID,
	)
	if err = tx.Commit(ctx); err != nil {
		return
	}
	return nil
}

// GetAnnotationIDForUnitTesting in only for unit testing purpose and will return annotation_id field
func GetAnnotationIDForUnitTesting() (z dbCommon, err error) {
	ctx := context.Background()
	db, err := pgxpool.Connect(ctx, commons.BuildDSN())
	if err != nil {
		return
	}
	defer db.Close()

	err = db.QueryRow(
		ctx,
		"SELECT * FROM annotations LIMIT 1",
	).Scan(
		&z.Annotation_id,
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
func (p *searchAnnotations) search() (z []dbList, err error) {
	ctx := context.Background()
	db, err := pgxpool.Connect(ctx, commons.BuildDSN())
	if err != nil {
		return
	}
	defer db.Close()

	rows, err := db.Query(
		ctx,
		"SELECT a.*, (SELECT COUNT(a.annotation_id) FROM annotations a LEFT JOIN projects p ON a.project_id = p.project_id WHERE a.key LIKE '%' || $1 || '%' OR a.value ILIKE '%' || $1 || '%' OR p.project_name ILIKE '%' || $1 || '%') total, p.project_name FROM annotations a LEFT JOIN projects p ON a.project_id = p.project_id WHERE a.key ILIKE '%' || $1 || '%' OR a.value ILIKE '%' || $1 || '%' OR p.project_name ILIKE '%' || $1 || '%' ORDER BY a.date DESC OFFSET $2 LIMIT $3",
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
			&x.Annotation_id,
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

// listByProjectID handle requirements to list annotations by project id
func (p *listAnnotationsByProjectID) listByProjectID() (z []dbRead, err error) {
	ctx := context.Background()
	db, err := pgxpool.Connect(ctx, commons.BuildDSN())
	if err != nil {
		return
	}
	defer db.Close()

	rows, err := db.Query(
		ctx,
		"SELECT a.*, p.project_name FROM annotations a LEFT JOIN projects p ON a.project_id = p.project_id WHERE a.project_id = $1 ORDER BY a.date DESC",
		p.ProjectID,
	)
	if err != nil && err.Error() != pgx.ErrNoRows.Error() {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var x dbRead
		if err = rows.Scan(
			&x.Annotation_id,
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
