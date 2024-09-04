-- Write your migrate up statements here
CREATE TABLE IF NOT EXISTS aluno (
    codigo SERIAL PRIMARY KEY,
    nome VARCHAR(50) NOT NULL
);

---- create above / drop below ----
DROP TABLE IF EXISTS aluno;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
