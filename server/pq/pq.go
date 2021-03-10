package pq

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func UsePostgres(ctx context.Context) bool {
	_, ok := os.LookupEnv("DB_USER")
	return ok
}

func Connect(ctx context.Context) (*sql.DB, error) {
	user, ok := os.LookupEnv("DB_USER")
	if !ok {
		return nil, errors.New("DB_USER not defined")
	}
	password, ok := os.LookupEnv("DB_PASS")
	if !ok {
		return nil, errors.New("DB_PASS not defined")
	}
	instance, ok := os.LookupEnv("DB_INSTANCE")
	if !ok {
		return nil, errors.New("DB_INSTANCE not defined")
	}
	name, ok := os.LookupEnv("DB_NAME")
	if !ok {
		return nil, errors.New("DB_NAME not defined")
	}
	socketDir, ok := os.LookupEnv("DB_SOCKET_DIR")
	if !ok {
		socketDir = "/cloudsql"
	}

	dbURI := fmt.Sprintf("user=%s password=%s database=%s host=%s/%s", user, password, name, socketDir, instance)

	dbPool, err := sql.Open("pgx", dbURI)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %v", err)
	}

	dbPool.SetMaxIdleConns(5)
	dbPool.SetMaxOpenConns(7)
	dbPool.SetConnMaxLifetime(1800)

	return dbPool, nil
}
