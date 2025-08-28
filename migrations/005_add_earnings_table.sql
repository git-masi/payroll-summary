CREATE TABLE IF NOT EXISTS earnings(
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    amount numeric(5, 4) NOT NULL,
    date_of_work date NOT NULL,
    payroll_id bigint REFERENCES payrolls(id) NOT NULL,
    worker_id bigint REFERENCES workers(id) NOT NULL,
    crew_id bigint REFERENCES crews(id),
    hours_worked numeric(2, 4),
    hours_offered numeric(2, 4),
    piece_units numeric(8, 4)
);

---- create above / drop below ----
DROP TABLE IF EXISTS earnings;

