package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreatePool(dsn string) (pool *pgxpool.Pool) {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatal("Failed to create database pool", err.Error())
	}
	return pool
}

func TestConnection(pool *pgxpool.Pool) {
	ctx := context.Background()
	_, err := pool.Acquire(ctx)
	if err != nil {
		log.Fatalf("failed to connect to postgres database: %s", err.Error())
	}
	fmt.Println("Connected to Postgres Database")
}
