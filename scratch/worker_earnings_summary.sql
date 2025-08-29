WITH cte_sum AS (
    SELECT
        worker_id,
        payroll_id,
        sum(amount) AS total_earnings,
        sum(hours_worked) AS total_hours_worked,
        sum(hours_offered) AS total_hours_offered,
        sum(piece_units) AS total_piece_units
    FROM
        earnings
        JOIN payrolls ON payrolls.id = earnings.payroll_id
    WHERE
        payrolls.status = 'paid'::payroll_status
    GROUP BY
        worker_id,
        payroll_id
)
SELECT
    worker_id,
    payroll_id,
    total_earnings,
    percent_rank() OVER (PARTITION BY payroll_id ORDER BY total_earnings DESC) AS earnings_percent_rank,
    cume_dist() OVER (PARTITION BY payroll_id ORDER BY total_earnings DESC) AS earnings_cume_dist,
    ntile(10) OVER (PARTITION BY payroll_id ORDER BY total_earnings DESC) AS decile_rank
FROM
    cte_sum;

-- ### How to read it:
--
-- * `percent_rank` → `1.0` means tied with the highest earner, `0.0` means lowest.
-- * `cume_dist` → fraction of workers with earnings **greater than or equal** to this worker.
-- * `decile_rank` → puts workers into 10 groups, `1 = top 10% earners`, `10 = bottom 10%`.
