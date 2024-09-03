// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: queries.sql

package pgstore

import (
	"context"
)

const getAllStudents = `-- name: GetAllStudents :many

select codigo, nome from aluno
`

func (q *Queries) GetAllStudents(ctx context.Context) ([]Aluno, error) {
	rows, err := q.db.Query(ctx, getAllStudents)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Aluno
	for rows.Next() {
		var i Aluno
		if err := rows.Scan(&i.Codigo, &i.Nome); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
