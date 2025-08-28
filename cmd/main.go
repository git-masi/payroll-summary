package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"payroll-summary/cmd/repo"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dsn := flag.String("dsn", "", "Postgres dsn")
	numWorkers := flag.Int("num_workers", 1000, "Number of workers to add")
	numCrews := flag.Int("num_crews", 20, "Number of workers to add")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))

	dbpool, err := connectDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer dbpool.Close()

	queries := repo.New(dbpool)

	err = createWorkers(queries, *numWorkers)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	logger.Info("Workers created", slog.Int("num_workers", *numWorkers))

	err = createCrews(queries, *numCrews)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	logger.Info("Crews created", slog.Int("num_crews", *numCrews))

	startTime := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
	endTime := startTime.AddDate(1, 0, 0)

	err = createMonthlyPayrolls(queries, startTime, endTime)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	logger.Info("Monthly payrolls created")

	logger.Info("Done")
}

func connectDB(dsn string) (*pgxpool.Pool, error) {
	// We don't actually need a pool but it's nice to have
	dbpool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = dbpool.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return dbpool, nil
}

func createWorkers(queries *repo.Queries, numWorkers int) error {
	newWorkers := make([]repo.CreateWorkersParams, numWorkers)
	for n := range numWorkers {
		newWorkers[n] = repo.CreateWorkersParams{
			FirstName: gofakeit.FirstName(),
			LastName:  gofakeit.LastName(),
		}
	}

	_, err := queries.CreateWorkers(context.TODO(), newWorkers)
	if err != nil {
		return err

	}

	return nil
}

func createCrews(queries *repo.Queries, numCrews int) error {
	newCrews := make([]string, numCrews)
	for n := range numCrews {
		newCrews[n] = gofakeit.AdjectiveDescriptive() + " " + gofakeit.NounCommon()
	}

	_, err := queries.CreateCrews(context.TODO(), newCrews)
	if err != nil {
		return err
	}

	return nil
}

func createMonthlyPayrolls(queries *repo.Queries, startTime time.Time, endTime time.Time) error {
	currentTime := startTime

	var newMonthlyPayrolls []repo.CreatePayrollsParams

	for currentTime.Before(endTime) {
		var start pgtype.Date
		var end pgtype.Date

		start.Scan(currentTime)
		end.Scan(currentTime.AddDate(0, 1, 0).AddDate(0, 0, -1))

		newMonthlyPayrolls = append(newMonthlyPayrolls, repo.CreatePayrollsParams{
			PayPeriod:   repo.PayrollPayPeriodMonthly,
			PeriodStart: start,
			PeriodEnd:   end,
		})

		currentTime = currentTime.AddDate(0, 1, 0)
	}

	_, err := queries.CreatePayrolls(context.TODO(), newMonthlyPayrolls)
	if err != nil {
		return err
	}

	return nil
}
