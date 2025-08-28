package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"payroll-summary/cmd/repo"
	"strconv"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dsn := flag.String("dsn", "", "Postgres dsn")
	numWorkers := flag.Int("num_workers", 1000, "Number of workers to add")
	numCrews := flag.Int("num_crews", 20, "Number of workers to add")
	shouldCreatePayrolls := flag.Bool("should_create_payrolls", true, "Whether or not to create new payrolls")
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

	if *shouldCreatePayrolls {
		startTime := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
		endTime := startTime.AddDate(1, 0, 0)

		err = createMonthlyPayrolls(queries, startTime, endTime)
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}
		logger.Info("Monthly payrolls created")

		err = createBiweeklyPayrolls(queries, startTime, endTime)
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}
		logger.Info("Biweekly payrolls created")

		err = createWeeklyPayrolls(queries, startTime, endTime)
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}
		logger.Info("Weekly payrolls created")
	}

	err = createEarnings(queries, logger)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	logger.Info("Earnings created")

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

func createBiweeklyPayrolls(queries *repo.Queries, startTime time.Time, endTime time.Time) error {
	currentTime := startTime

	var newBiweeklyPayrolls []repo.CreatePayrollsParams

	for currentTime.Before(endTime) {
		twoWeeks := currentTime.AddDate(0, 0, 14)

		var start pgtype.Date
		var end pgtype.Date

		start.Scan(currentTime)
		end.Scan(twoWeeks)

		newBiweeklyPayrolls = append(newBiweeklyPayrolls, repo.CreatePayrollsParams{
			PayPeriod:   repo.PayrollPayPeriodBiweekly,
			PeriodStart: start,
			PeriodEnd:   end,
		})

		currentTime = twoWeeks
	}

	_, err := queries.CreatePayrolls(context.TODO(), newBiweeklyPayrolls)
	if err != nil {
		return err
	}

	return nil
}

func createWeeklyPayrolls(queries *repo.Queries, startTime time.Time, endTime time.Time) error {
	currentTime := startTime

	var newWeeklyPayrolls []repo.CreatePayrollsParams

	for currentTime.Before(endTime) {
		oneWeek := currentTime.AddDate(0, 0, 7)

		var start pgtype.Date
		var end pgtype.Date

		start.Scan(currentTime)
		end.Scan(oneWeek)

		newWeeklyPayrolls = append(newWeeklyPayrolls, repo.CreatePayrollsParams{
			PayPeriod:   repo.PayrollPayPeriodWeekly,
			PeriodStart: start,
			PeriodEnd:   end,
		})

		currentTime = oneWeek
	}

	_, err := queries.CreatePayrolls(context.TODO(), newWeeklyPayrolls)
	if err != nil {
		return err
	}

	return nil
}

func createEarnings(queries *repo.Queries, logger *slog.Logger) error {
	workerIDs, _ := queries.GetWorkerIDs(context.TODO())
	crewIDs, _ := queries.GetCrewIDs(context.TODO())
	payrolls, _ := queries.GetPayrolls(context.TODO())

	numPayrolls := len(payrolls)

	for i, p := range payrolls {
		currentTime := p.PeriodStart.Time

		newEarnings := []repo.CreateEarningsParams{}

		// For each date in the payroll
		for currentTime.Before(p.PeriodEnd.Time) {
			crewID := crewIDs[rand.Intn(len(crewIDs))]

			// Create earnings for each worker
			for _, workerID := range workerIDs {
				params, err := createEarningParams(currentTime, p.ID, workerID, crewID)
				if err != nil {
					return err
				}
				newEarnings = append(newEarnings, *params)
			}

			currentTime = currentTime.AddDate(0, 0, 1)
		}

		_, err := queries.CreateEarnings(context.TODO(), newEarnings)
		if err != nil {
			return err
		}

		logger.Info(fmt.Sprintf("created earnings for payroll number %d of %d", i+1, numPayrolls), slog.Int("num_earnings", len(newEarnings)))
	}

	return nil
}

func createEarningParams(currentTime time.Time, payrollID int64, workerID int64, crewID int64) (*repo.CreateEarningsParams, error) {
	pieceWork := gofakeit.FlipACoin() == "Heads"

	var amount pgtype.Numeric
	err := amount.Scan(strconv.FormatFloat(gofakeit.Price(10, 500), 'f', 4, 32))
	if err != nil {
		return nil, err
	}

	var dateOfWork pgtype.Date
	err = dateOfWork.Scan(currentTime)
	if err != nil {
		return nil, err
	}

	var hoursWorked pgtype.Numeric
	if !pieceWork {
		err = hoursWorked.Scan(strconv.FormatFloat(gofakeit.Float64Range(4, 12), 'f', 4, 32))
		if err != nil {
			return nil, err
		}
	}

	var hoursOffered pgtype.Numeric
	if !pieceWork {
		err = hoursOffered.Scan(strconv.FormatFloat(gofakeit.Float64Range(4, 12), 'f', 4, 32))
		if err != nil {
			return nil, err
		}
	}

	var pieceUnits pgtype.Numeric
	if pieceWork {
		err = pieceUnits.Scan(strconv.FormatFloat(gofakeit.Float64Range(100, 1000), 'f', 4, 32))
		if err != nil {
			return nil, err
		}
	}

	var crew pgtype.Int8
	if pieceWork {
		crew.Scan(crewID)
	}

	return &repo.CreateEarningsParams{
		Amount:       amount,
		DateOfWork:   dateOfWork,
		PayrollID:    payrollID,
		WorkerID:     workerID,
		CrewID:       crew,
		HoursWorked:  hoursWorked,
		HoursOffered: hoursOffered,
		PieceUnits:   pieceUnits,
	}, nil
}
