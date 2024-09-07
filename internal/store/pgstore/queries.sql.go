// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: queries.sql

package pgstore

import (
	"context"
)

const alunosMatriculados = `-- name: AlunosMatriculados :one
SELECT 
    COUNT(ca.codigo_aluno) AS numero_alunos
FROM 
    curso_aluno ca
WHERE 
    ca.codigo_curso = $1
`

func (q *Queries) AlunosMatriculados(ctx context.Context, codigoCurso int32) (int64, error) {
	row := q.db.QueryRow(ctx, alunosMatriculados, codigoCurso)
	var numero_alunos int64
	err := row.Scan(&numero_alunos)
	return numero_alunos, err
}

const createAluno = `-- name: CreateAluno :one
insert into aluno
    ("nome") 
 values 
    ($1) 
 RETURNING "codigo"
`

func (q *Queries) CreateAluno(ctx context.Context, nome string) (int32, error) {
	row := q.db.QueryRow(ctx, createAluno, nome)
	var codigo int32
	err := row.Scan(&codigo)
	return codigo, err
}

const createCurso = `-- name: CreateCurso :one
insert into curso 
    ("descricao", "ementa")
values 
    ($1, $2) RETURNING codigo, descricao, ementa
`

type CreateCursoParams struct {
	Descricao string
	Ementa    string
}

func (q *Queries) CreateCurso(ctx context.Context, arg CreateCursoParams) (Curso, error) {
	row := q.db.QueryRow(ctx, createCurso, arg.Descricao, arg.Ementa)
	var i Curso
	err := row.Scan(&i.Codigo, &i.Descricao, &i.Ementa)
	return i, err
}

const cursosMatriculados = `-- name: CursosMatriculados :many
SELECT 
    c.codigo
FROM 
    curso c 
LEFT JOIN 
    curso_aluno ca ON c.codigo = ca.codigo_curso
WHERE 
    ca.codigo_aluno = $1
GROUP BY 
    c.codigo
`

func (q *Queries) CursosMatriculados(ctx context.Context, codigoAluno int32) ([]int32, error) {
	rows, err := q.db.Query(ctx, cursosMatriculados, codigoAluno)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []int32
	for rows.Next() {
		var codigo int32
		if err := rows.Scan(&codigo); err != nil {
			return nil, err
		}
		items = append(items, codigo)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const deleteAluno = `-- name: DeleteAluno :exec
DELETE FROM 
    aluno 
WHERE 
    codigo = $1
`

func (q *Queries) DeleteAluno(ctx context.Context, codigo int32) error {
	_, err := q.db.Exec(ctx, deleteAluno, codigo)
	return err
}

const deleteCurso = `-- name: DeleteCurso :exec
DELETE FROM 
    curso 
WHERE 
    codigo=$1
`

func (q *Queries) DeleteCurso(ctx context.Context, codigo int32) error {
	_, err := q.db.Exec(ctx, deleteCurso, codigo)
	return err
}

const getAllAlunos = `-- name: GetAllAlunos :many
select codigo, nome from aluno
`

func (q *Queries) GetAllAlunos(ctx context.Context) ([]Aluno, error) {
	rows, err := q.db.Query(ctx, getAllAlunos)
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

const getAluno = `-- name: GetAluno :one
SELECT codigo, nome FROM 
    aluno a
WHERE 
    a.codigo = $1
`

func (q *Queries) GetAluno(ctx context.Context, codigo int32) (Aluno, error) {
	row := q.db.QueryRow(ctx, getAluno, codigo)
	var i Aluno
	err := row.Scan(&i.Codigo, &i.Nome)
	return i, err
}

const getCurso = `-- name: GetCurso :one
SELECT 
    c.codigo, c.descricao, c.ementa,
    CAST(COALESCE(COUNT(ca.codigo_aluno), 0) AS INTEGER) AS matriculas 
FROM
    curso c
LEFT JOIN 
    curso_aluno ca ON c.codigo = ca.codigo_curso
WHERE 
    c.codigo = $1
GROUP BY 
    c.codigo
`

type GetCursoRow struct {
	Codigo     int32
	Descricao  string
	Ementa     string
	Matriculas int32
}

func (q *Queries) GetCurso(ctx context.Context, codigo int32) (GetCursoRow, error) {
	row := q.db.QueryRow(ctx, getCurso, codigo)
	var i GetCursoRow
	err := row.Scan(
		&i.Codigo,
		&i.Descricao,
		&i.Ementa,
		&i.Matriculas,
	)
	return i, err
}

const getCursoByCodigo = `-- name: GetCursoByCodigo :one
select codigo, descricao, ementa from curso 
where 
    codigo = $1
`

func (q *Queries) GetCursoByCodigo(ctx context.Context, codigo int32) (Curso, error) {
	row := q.db.QueryRow(ctx, getCursoByCodigo, codigo)
	var i Curso
	err := row.Scan(&i.Codigo, &i.Descricao, &i.Ementa)
	return i, err
}

const getCursos = `-- name: GetCursos :many
SELECT 
    c.codigo, c.descricao, c.ementa, 
    CAST(COALESCE(COUNT(ca.codigo_aluno), 0) AS INTEGER) AS matriculas
FROM 
    curso c
LEFT JOIN 
    curso_aluno ca ON c.codigo = ca.codigo_curso
GROUP BY 
    c.codigo
`

type GetCursosRow struct {
	Codigo     int32
	Descricao  string
	Ementa     string
	Matriculas int32
}

func (q *Queries) GetCursos(ctx context.Context) ([]GetCursosRow, error) {
	rows, err := q.db.Query(ctx, getCursos)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetCursosRow
	for rows.Next() {
		var i GetCursosRow
		if err := rows.Scan(
			&i.Codigo,
			&i.Descricao,
			&i.Ementa,
			&i.Matriculas,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getDashBoardInfo = `-- name: GetDashBoardInfo :one
WITH total_alunos AS (
    SELECT COALESCE(COUNT(*), 0) AS total FROM aluno
),
total_cursos AS (
    SELECT COALESCE(COUNT(*), 0) AS total FROM curso
),
total_matriculas AS (
    SELECT COALESCE(COUNT(*), 0) AS total FROM curso_aluno
),
matriculas_por_curso AS (
    SELECT 
        c.descricao AS curso, 
        COALESCE(COUNT(ca.codigo_aluno), 0) AS total_matriculas
    FROM 
        curso c
    LEFT JOIN 
        curso_aluno ca ON c.codigo = ca.codigo_curso
    GROUP BY c.descricao
),
alunos_com_matricula AS (
    SELECT DISTINCT
        a.nome, 
        a.codigo 
    FROM 
        aluno a
    INNER JOIN 
        curso_aluno ca ON a.codigo = ca.codigo_aluno
),
alunos_sem_matricula AS (
    SELECT DISTINCT
        a.nome, 
        a.codigo 
    FROM 
        aluno a
    LEFT JOIN 
        curso_aluno ca ON a.codigo = ca.codigo_aluno
    WHERE 
        ca.codigo_aluno IS NULL
)
SELECT 
    (SELECT total FROM total_alunos) AS total_alunos,
    (SELECT total FROM total_cursos) AS total_cursos,
    (SELECT total FROM total_matriculas) AS total_matriculas,
    COALESCE((SELECT json_agg(matriculas_por_curso) FROM matriculas_por_curso), '[]'::json) AS matriculas_por_curso,
    COALESCE((SELECT json_agg(alunos_com_matricula) FROM alunos_com_matricula), '[]'::json) AS alunos_com_matricula,
    COALESCE((SELECT json_agg(alunos_sem_matricula) FROM alunos_sem_matricula), '[]'::json) AS alunos_sem_matricula
`

type GetDashBoardInfoRow struct {
	TotalAlunos        interface{}
	TotalCursos        interface{}
	TotalMatriculas    interface{}
	MatriculasPorCurso interface{}
	AlunosComMatricula interface{}
	AlunosSemMatricula interface{}
}

func (q *Queries) GetDashBoardInfo(ctx context.Context) (GetDashBoardInfoRow, error) {
	row := q.db.QueryRow(ctx, getDashBoardInfo)
	var i GetDashBoardInfoRow
	err := row.Scan(
		&i.TotalAlunos,
		&i.TotalCursos,
		&i.TotalMatriculas,
		&i.MatriculasPorCurso,
		&i.AlunosComMatricula,
		&i.AlunosSemMatricula,
	)
	return i, err
}

const matricularAluno = `-- name: MatricularAluno :one
insert into curso_aluno 
    ("codigo_aluno", "codigo_curso")
values
    ($1, $2)
RETURNING
    "codigo"
`

type MatricularAlunoParams struct {
	CodigoAluno int32
	CodigoCurso int32
}

func (q *Queries) MatricularAluno(ctx context.Context, arg MatricularAlunoParams) (int32, error) {
	row := q.db.QueryRow(ctx, matricularAluno, arg.CodigoAluno, arg.CodigoCurso)
	var codigo int32
	err := row.Scan(&codigo)
	return codigo, err
}

const matriculasPorAluno = `-- name: MatriculasPorAluno :one
SELECT 
    COUNT(ca.codigo_curso) AS numero_cursos
FROM 
    curso_aluno ca
WHERE 
    ca.codigo_aluno = $1
`

func (q *Queries) MatriculasPorAluno(ctx context.Context, codigoAluno int32) (int64, error) {
	row := q.db.QueryRow(ctx, matriculasPorAluno, codigoAluno)
	var numero_cursos int64
	err := row.Scan(&numero_cursos)
	return numero_cursos, err
}

const searchAlunos = `-- name: SearchAlunos :many
SELECT 
    nome 
FROM 
    aluno 
WHERE 
    unaccent(nome) ILIKE '%' || unaccent($1) || '%'
    OR to_tsvector('simple', unaccent(nome)) @@ to_tsquery('simple', unaccent($1))
`

func (q *Queries) SearchAlunos(ctx context.Context, unaccent string) ([]string, error) {
	rows, err := q.db.Query(ctx, searchAlunos, unaccent)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var nome string
		if err := rows.Scan(&nome); err != nil {
			return nil, err
		}
		items = append(items, nome)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateCurso = `-- name: UpdateCurso :one
update curso 
set 
    descricao=$2,
    ementa=$3
where
    codigo=$1
RETURNING codigo, descricao, ementa
`

type UpdateCursoParams struct {
	Codigo    int32
	Descricao string
	Ementa    string
}

func (q *Queries) UpdateCurso(ctx context.Context, arg UpdateCursoParams) (Curso, error) {
	row := q.db.QueryRow(ctx, updateCurso, arg.Codigo, arg.Descricao, arg.Ementa)
	var i Curso
	err := row.Scan(&i.Codigo, &i.Descricao, &i.Ementa)
	return i, err
}

const updateNomeAluno = `-- name: UpdateNomeAluno :one

update aluno 
set 
    nome=$2
where 
    codigo=$1
RETURNING codigo, nome
`

type UpdateNomeAlunoParams struct {
	Codigo int32
	Nome   string
}

func (q *Queries) UpdateNomeAluno(ctx context.Context, arg UpdateNomeAlunoParams) (Aluno, error) {
	row := q.db.QueryRow(ctx, updateNomeAluno, arg.Codigo, arg.Nome)
	var i Aluno
	err := row.Scan(&i.Codigo, &i.Nome)
	return i, err
}
