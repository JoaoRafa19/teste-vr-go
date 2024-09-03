-- Write your migrate up statements here
CREATE TABLE IF NOT EXISTS curso (
    "codigo"    SERIAL       PRIMARY KEY            ,
    "descricao" VARCHAR(50)                 NOT NULL,
    "ementa"    TEXT                        NOT NULL
);
---- create above / drop below ----
DROP TABLE curso;
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
