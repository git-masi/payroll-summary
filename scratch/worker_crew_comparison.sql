WITH cte_sum AS (
    SELECT
        worker_id,
        crew_id,
        sum(amount) AS total_earnings,
        sum(hours_worked) AS total_hours_worked,
        sum(hours_offered) AS total_hours_offered,
        sum(piece_units) AS total_piece_units
    FROM
        earnings
        JOIN payrolls ON payrolls.id = earnings.payroll_id
    WHERE
        payroll_id = 42
        AND crew_id IS NOT NULL
    GROUP BY
        worker_id,
        crew_id
)
SELECT
    worker_id,
    crew_id,
    total_piece_units,
    percent_rank() OVER (PARTITION BY crew_id ORDER BY total_piece_units DESC) AS piece_units_percent_rank,
    cume_dist() OVER (PARTITION BY crew_id ORDER BY total_piece_units DESC) AS piece_units_cume_dist,
    ntile(10) OVER (PARTITION BY crew_id ORDER BY total_piece_units DESC) AS decile_rank
FROM
    cte_sum
ORDER BY
    crew_id,
    decile_rank;

