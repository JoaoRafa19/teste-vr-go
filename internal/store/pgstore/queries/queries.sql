-- name: GetAllAlunos :many
select * from aluno;

-- name: CreateAluno :one 
insert into aluno
    ("nome") 
 values 
    ('$1') 
 RETURNING "id";

-- name: UpdateNomeAluno :one

update aluno 
set 
    nome=$2
where 
    codigo=$1
RETURNING "codigo";


-- name: CreateCurso :one
insert into curso 
    ("descricao", "ementa")
values 
    ($1, $2)
RETURNING "codigo";

-- name: GetCursos :many
select * from curso;

-- UpdateCurso :one 
update curso 
set 
    descricao=$2,
    ementa=$3
where
    id=$1
RETURNING *;


