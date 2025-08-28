---- tern: disable-tx ----
CREATE INDEX CONCURRENTLY IF NOT EXISTS payroll_period_start_idx ON payrolls(period_start);

---- create above / drop below ----
DROP INDEX IF EXISTS payroll_period_start_idx;

