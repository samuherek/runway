package db

import (
	"context"
	"database/sql"
	"log"
	"os"
	"runway/db/dbgen"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func getDatabaseUrl() string {
	env, ok := os.LookupEnv("APP_DATABASE_URL")
	if !ok {
		log.Fatal("ERROR: Failed to load APP_DATABASE_URL.")
	}
	return env
}

func openConnection(url string) *sql.DB {
	db, err := sql.Open("pgx", url)
	if err != nil {
		log.Fatalf("Failed to connect DB: %v", err)
	}

	return db
}

type DbService struct {
	pool    *sql.DB
	Queries *dbgen.Queries
}

func NewDbService() *DbService {
	url := getDatabaseUrl()
	db := openConnection(url)

	// TODO: Run migrations

	queries := dbgen.New(db)

	return &DbService{
		pool:    db,
		Queries: queries,
	}
}

func (s *DbService) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return s.pool.BeginTx(ctx, nil)
}

func (s *DbService) Close() error {
	return s.pool.Close()
}
