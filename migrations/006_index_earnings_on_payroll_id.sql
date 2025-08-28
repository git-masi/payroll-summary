---- tern: disable-tx ----
CREATE INDEX CONCURRENTLY IF NOT EXISTS earnings_payroll_id_idx ON earnings(payroll_id);

---- create above / drop below ----
DROP INDEX IF EXISTS earnings_payroll_id_idx;

