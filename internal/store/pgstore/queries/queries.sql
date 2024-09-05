-- name: GetAllAlunos :many
select * from aluno;

-- name: GetAluno :one
SELECT * FROM 
    aluno a
WHERE 
    a.codigo = $1;
    
-- name: DeleteAluno :exec
DELETE FROM 
    aluno 
WHERE 
    codigo = $1;

-- name: DeleteCurso :exec
DELETE FROM 
    curso 
WHERE 
    codigo=$1;


-- name: CreateAluno :one 
insert into aluno
    ("nome") 
 values 
    ($1) 
 RETURNING "codigo";

-- name: UpdateNomeAluno :one

update aluno 
set 
    nome=$2
where 
    codigo=$1
RETURNING *;


-- name: CreateCurso :one
insert into curso 
    ("descricao", "ementa")
values 
    ($1, $2) RETURNING *;



-- name: UpdateCurso :one 
update curso 
set 
    descricao=$2,
    ementa=$3
where
    codigo=$1
RETURNING *;

-- name: MatricularAluno :one 
insert into curso_aluno 
    ("codigo_aluno", "codigo_curso")
values
    ($1, $2)
RETURNING
    "codigo";


-- name: AlunosMatriculados :one
SELECT 
    COUNT(ca.codigo_aluno) AS numero_alunos
FROM 
    curso_aluno ca
WHERE 
    ca.codigo_curso = $1;

-- name: MatriculasPorAluno :one
SELECT 
    COUNT(ca.codigo_curso) AS numero_cursos
FROM 
    curso_aluno ca
WHERE 
    ca.codigo_aluno = $1;


-- name: CursosMatriculados :many
SELECT 
    c.codigo
FROM 
    curso c 
LEFT JOIN 
    curso_aluno ca ON c.codigo = ca.codigo_curso
WHERE 
    ca.codigo_aluno = $1
GROUP BY 
    c.codigo;
    

-- name: GetCursos :many
SELECT 
    c.*, 
    CAST(COALESCE(COUNT(ca.codigo_aluno), 0) AS INTEGER) AS matriculas
FROM 
    curso c
LEFT JOIN 
    curso_aluno ca ON c.codigo = ca.codigo_curso
GROUP BY 
    c.codigo;

-- name: GetCurso :one
SELECT 
    c.*,
    CAST(COALESCE(COUNT(ca.codigo_aluno), 0) AS INTEGER) AS matriculas 
FROM
    curso c
LEFT JOIN 
    curso_aluno ca ON c.codigo = ca.codigo_curso
WHERE 
    c.codigo = $1
GROUP BY 
    c.codigo;

-- name: SearchAlunos :many
SELECT 
    nome 
FROM 
    aluno 
WHERE 
    unaccent(nome) ILIKE '%' || unaccent($1) || '%'
    OR to_tsvector('simple', unaccent(nome)) @@ to_tsquery('simple', unaccent($1));
