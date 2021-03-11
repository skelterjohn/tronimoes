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
