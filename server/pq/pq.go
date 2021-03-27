package pq

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"

	"cloud.google.com/go/compute/metadata"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/skelterjohn/tronimoes/server/util"
)

func getServiceAccountAccessToken(ctx context.Context) (string, error) {
	return metadata.Get("/instance/service-accounts/default/token")
}
func getServiceAccountEmail(ctx context.Context) (string, error) {
	return metadata.Get("/instance/service-accounts/default/email")
}

func UsePostgres(ctx context.Context) bool {
	_, ok := os.LookupEnv("DB_INSTANCE")
	return ok
}

func Connect(ctx context.Context) (*sql.DB, error) {
	user, ok := os.LookupEnv("DB_USER")
	if !ok {
		serviceAccountEmail, err := getServiceAccountEmail(ctx)
		if err != nil {
			return nil, util.Annotate(err, "DB_USER not defined and could not use metadata")
		}
		user = serviceAccountEmail
	}
	password, ok := os.LookupEnv("DB_PASS")
	if !ok {
		token, err := getServiceAccountAccessToken(ctx)
		if err != nil {
			return nil, util.Annotate(err, "DB_PASS not defined and could not use metadata")
		}
		password = token
	}
	instance, instance_ok := os.LookupEnv("DB_INSTANCE")
	host, host_ok := os.LookupEnv("DB_HOST")
	if !instance_ok && !host_ok {
		return nil, errors.New("neither DB_INSTANCE nor DB_HOST defined")
	}

	name, ok := os.LookupEnv("DB_NAME")
	if !ok {
		return nil, errors.New("DB_NAME not defined")
	}
	socketDir, ok := os.LookupEnv("DB_SOCKET_DIR")
	if !ok {
		socketDir = "/cloudsql"
	}

	var dbURI string
	if host_ok {
		dbURI = fmt.Sprintf("user=%s password=%s database=%s host=%s port=5432", user, password, name, host)
		fmt.Printf("Connecting to database %s as %s @ %s\n", name, user, host)

	}
	if instance_ok {
		dbURI = fmt.Sprintf("user=%s password=%s database=%s host=%s/%s", user, password, name, socketDir, instance)
		fmt.Printf("Connecting to database %s as %s @ %s/%s\n", name, user, socketDir, instance)
	}

	dbPool, err := sql.Open("pgx", dbURI)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %v", err)
	}

	dbPool.SetMaxIdleConns(5)
	dbPool.SetMaxOpenConns(7)
	dbPool.SetConnMaxLifetime(1800)

	return dbPool, nil
}
