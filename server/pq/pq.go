package pq

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"

	"cloud.google.com/go/compute/metadata"
	secretsv1 "cloud.google.com/go/secretmanager/apiv1"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/skelterjohn/tronimoes/server/util"
	secretpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

func getServiceAccountAccessToken(ctx context.Context) (string, error) {
	return metadata.Get("/instance/service-accounts/default/token")
}
func getServiceAccountEmail(ctx context.Context) (string, error) {
	return metadata.Get("/instance/service-accounts/default/email")
}

func UsePostgres(ctx context.Context) bool {
	fmt.Printf("DB_INSTANCE=%s\n", os.Getenv("DB_INSTANCE"))
	fmt.Printf("DB_HOST=%s\n", os.Getenv("DB_HOST"))
	if _, ok := os.LookupEnv("DB_INSTANCE"); ok {
		return true
	}
	if _, ok := os.LookupEnv("DB_HOST"); ok {
		return true
	}
	return false
}

func Connect(ctx context.Context) (*sql.DB, error) {
	user, user_ok := os.LookupEnv("DB_USER")
	if !user_ok {
		serviceAccountEmail, err := getServiceAccountEmail(ctx)
		if err != nil {
			return nil, util.Annotate(err, "DB_USER not defined and could not use metadata")
		}
		user = serviceAccountEmail
	}

	passwordSecret, secret_ok := os.LookupEnv("DB_PASS_SECRET")
	password, password_ok := os.LookupEnv("DB_PASS")
	if secret_ok && password_ok {
		return nil, errors.New("must have at most one of DB_PASS and DB_PASS_SECRET defined")
	}
	if !user_ok && (secret_ok || password_ok) {
		return nil, errors.New("may not define DB_PASS or DB_PASS_SECRET without DB_USER")
	}
	if secret_ok {
		c, err := secretsv1.NewClient(ctx)
		if err != nil {
			return nil, util.Annotate(err, "could not get secrets client")
		}
		s, err := c.AccessSecretVersion(ctx, &secretpb.AccessSecretVersionRequest{
			Name: passwordSecret,
		})
		if err != nil {
			return nil, util.Annotate(err, "could not fetch secret")
		}
		password = string(s.GetPayload().GetData())
	}

	if !user_ok {
		token, err := getServiceAccountAccessToken(ctx)
		if err != nil {
			return nil, util.Annotate(err, "DB_PASS not defined and could not use metadata")
		}
		password = token
	}
	instance, instance_ok := os.LookupEnv("DB_INSTANCE")
	host, host_ok := os.LookupEnv("DB_HOST")
	if instance_ok == host_ok {
		return nil, errors.New("must have exactly one of DB_INSTANCE or DB_HOST defined")
	}

	name, ok := os.LookupEnv("DB_NAME")
	if !ok {
		return nil, errors.New("DB_NAME not defined")
	}
	socketDir, ok := os.LookupEnv("DB_SOCKET_DIR")
	if !ok {
		socketDir = "/cloudsql"
	}

	port, ok := os.LookupEnv("DB_PORT")
	if !ok {
		port = "5432"
	}

	var dbURI string
	if host_ok {
		dbURI = fmt.Sprintf("user=%s password=%s database=%s host=%s port=%s", user, password, name, host, port)
		fmt.Printf("Connecting to database %s as %s @ %s:%s\n", name, user, host, port)

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
