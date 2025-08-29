BEGIN;
CREATE temp TABLE payroll_summaries(
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    payroll_id bigint NOT NULL,
    total_earnings numeric(14, 4) NOT NULL,
    average_earnings numeric(14, 4) NOT NULL,
    median_earnings numeric(14, 4 NOT NULL)
) ON COMMIT DROP;
SAVEPOINT after_payroll_summaries;
WITH cte_sum AS (
    SELECT
        payroll_id,
        sum(amount) AS total_earnings
    FROM
        earnings
        JOIN payrolls ON payrolls.id = earnings.payroll_id
    WHERE
        payrolls.status = 'paid'::payroll_status
    GROUP BY
        payroll_id
),
cte_stats AS (
    SELECT
        avg(total_earnings) AS average_earnings,
        percentile_cont(0.5) WITHIN GROUP (ORDER BY total_earnings) AS median_earnings
    FROM
        cte_sum)
    INSERT INTO payroll_summaries(payroll_id, total_earnings, average_earnings, median_earnings)
    SELECT
        cte_sum.payroll_id,
        cte_sum.total_earnings,
        cte_stats.average_earnings,
        cte_stats.median_earnings
    FROM
        cte_sum
    CROSS JOIN cte_stats;
-- ROLLBACK TO SAVEPOINT after_payroll_summaries;
SELECT
    *
FROM
    payroll_summaries;
COMMIT;

