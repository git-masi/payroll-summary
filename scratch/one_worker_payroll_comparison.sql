-- update payrolls set status = 'draft'::payroll_status where id = 12;
WITH cte_current_sum AS (
    SELECT
        e.worker_id,
        e.payroll_id,
        p.pay_period,
        sum(e.amount) AS total_earnings
    FROM
        earnings e
        JOIN payrolls p ON p.id = e.payroll_id
    WHERE
        e.payroll_id = 12
        AND p.status = 'draft'::payroll_status
        AND p.pay_period = 'monthly'::payroll_pay_period
    GROUP BY
        e.worker_id,
        e.payroll_id,
        p.pay_period
),
cte_historical_sum AS (
    SELECT
        e.worker_id,
        e.payroll_id,
        p.pay_period,
        sum(e.amount) AS total_earnings
    FROM
        earnings e
        JOIN payrolls p ON p.id = e.payroll_id
    WHERE
        p.status = 'paid'::payroll_status
        AND p.pay_period = 'monthly'::payroll_pay_period
        AND e.worker_id IN (
            SELECT
                worker_id
            FROM
                cte_current_sum)
        GROUP BY
            e.worker_id,
            e.payroll_id,
            p.pay_period
),
cte_historical_stats AS (
    SELECT
        worker_id,
        avg(total_earnings) average_earnings
    FROM
        cte_historical_sum
    GROUP BY
        worker_id
)
SELECT
    c.worker_id,
    c.payroll_id,
    c.pay_period,
    c.total_earnings,
    s.average_earnings,
    c.total_earnings - s.average_earnings AS earnings_diff,
((c.total_earnings - s.average_earnings) / s.average_earnings) * 100 AS earnings_diff_percent
FROM
    cte_current_sum c
    JOIN cte_historical_stats s ON c.worker_id = s.worker_id
