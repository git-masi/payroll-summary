package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"payroll-summary/cmd/repo"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))

	dsn := flag.String("dsn", "", "postgresql dsn")
	flag.Parse()

	dbpool, err := connectDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer dbpool.Close()

	queries := repo.New(dbpool)

	numWorkers := 1
	newWorkers := make([]repo.CreateWorkersParams, numWorkers)
	for n := range numWorkers {
		newWorkers[n] = repo.CreateWorkersParams{
			FirstName: gofakeit.FirstName(),
			LastName:  gofakeit.LastName(),
		}
	}

	_, err = queries.CreateWorkers(context.TODO(), newWorkers)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	logger.Info("Workers created")

	numCrews := 20
	newCrews := make([]string, numCrews)
	for n := range numCrews {
		newCrews[n] = gofakeit.AdjectiveDescriptive() + " " + gofakeit.NounCommon()
	}

	_, err = queries.CreateCrews(context.TODO(), newCrews)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	logger.Info("Crews created")

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
