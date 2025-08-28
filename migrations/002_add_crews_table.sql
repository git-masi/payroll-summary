CREATE TABLE IF NOT EXISTS crews(
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name text NOT NULL
);

---- create above / drop below ----
DROP TABLE IF EXISTS crews;

