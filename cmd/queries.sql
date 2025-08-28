-- name: CreateWorkers :copyfrom
INSERT INTO workers(first_name, last_name)
    VALUES ($1, $2);

-- name: CreateCrews :copyfrom
INSERT INTO crews(name)
    VALUES ($1);

-- name: CreatePayrolls :copyfrom
INSERT INTO payrolls(pay_period, period_start, period_end)
    VALUES ($1, $2, $3);

-- name: CreateEarnings :copyfrom
INSERT INTO earnings(amount, date_of_work, payroll_id, worker_id, crew_id)
    VALUES ($1, $2, $3, $4, $5);

