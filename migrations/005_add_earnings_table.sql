CREATE TABLE IF NOT EXISTS earnings(
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    amount numeric(5, 4) NOT NULL,
    date_of_work date NOT NULL,
    payroll_id bigint REFERENCES payrolls(id) NOT NULL,
    worker_id bigint REFERENCES workers(id) NOT NULL,
    crew_id bigint REFERENCES crews(id)
);

---- create above / drop below ----
DROP TABLE IF EXISTS earnings;

