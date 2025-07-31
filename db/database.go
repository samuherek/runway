package db

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"os"
	"runway/db/dbgen"
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
	defer db.Close()

	return db
}

type DbService struct {
	pool    *sql.DB
	queries *dbgen.Queries
}

func NewDbService() *DbService {
	url := getDatabaseUrl()
	db := openConnection(url)

	// TODO: Run migrations

	queries := dbgen.New(db)

	return &DbService{
		pool:    db,
		queries: queries,
	}
}

func (dbs *DbService) Close() error {
	return dbs.pool.Close()
}
