-- Write your migrate up statements here
CREATE TABLE IF NOT EXISTS curso_aluno (
    "codigo" SERIAL PRIMARY KEY,
    "codigo_aluno" INTEGER NOT NULL,
    "codigo_curso" INTEGER NOT NULL,

    FOREIGN KEY (codigo_aluno) REFERENCES aluno(codigo) ON DELETE CASCADE ,
    FOREIGN KEY (codigo_curso) REFERENCES curso(codigo) ON DELETE CASCADE
);
---- create above / drop below ----
DROP TABLE IF EXISTS curso_aluno;
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
