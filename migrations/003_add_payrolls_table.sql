CREATE TYPE payroll_status AS enum(
    'draft',
    'pending',
    'paid',
    'void'
);

CREATE TYPE payroll_pay_period AS enum(
    'weekly',
    'biweekly',
    'monthly'
);

CREATE TABLE IF NOT EXISTS payrolls(
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    pay_period payroll_pay_period NOT NULL,
    period_start date NOT NULL,
    period_end date NOT NULL,
    status payroll_status NOT NULL DEFAULT 'draft' ::payroll_status
);

---- create above / drop below ----
DROP TABLE IF EXISTS payrolls;

DROP TYPE payroll_status;

DROP TYPE payroll_pay_period;

