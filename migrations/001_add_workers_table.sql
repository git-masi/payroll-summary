CREATE TABLE IF NOT EXISTS workers(
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    first_name text NOT NULL,
    last_name text NOT NULL
);

---- create above / drop below ----
DROP TABLE IF EXISTS workers;

