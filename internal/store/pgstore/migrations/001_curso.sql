-- Write your migrate up statements here
CREATE TABLE IF NOT EXISTS curso (
    "id"        uuid            PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    "descricao" VARCHAR(50)                 NOT NULL,
    "ementa"    TEXT                        NOT NULL
);
---- create above / drop below ----
DROP TABLE curso;
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
