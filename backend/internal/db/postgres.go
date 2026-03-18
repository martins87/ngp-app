package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool() (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	// pool, err := pgxpool.New(ctx, dsn)
	// if err != nil {
	// 	return nil, err
	// }

	// if err := pool.Ping(ctx); err != nil {
	// 	return nil, err
	// }

	// return pool, nil

	var pool *pgxpool.Pool
	var err error

	for range 15 {
		pool, err = pgxpool.New(ctx, dsn)
		if err == nil {
			err = pool.Ping(ctx)
			if err == nil {
				log.Println("Connected to database")
				return pool, nil
			}
		}

		log.Println("Database not ready yet, retrying...")
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("Could not connect to database after %w retries", err)
}