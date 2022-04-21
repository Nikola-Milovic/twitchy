package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type PgxIface interface {
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	Close(context.Context) error
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
}

func InitDb(ctx context.Context) (PgxIface, func() error, error) {
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"))

	fmt.Println(os.Getenv("POSTGRES_USER"))
	fmt.Println(os.Getenv("POSTGRES_PASSWORD"))
	fmt.Println(os.Getenv("POSTGRES_HOST"))
	fmt.Println(os.Getenv("POSTGRES_PORT"))
	fmt.Println(os.Getenv("POSTGRES_DB"))
	fmt.Println(dbUrl)

	conn, err := pgx.Connect(context.Background(), dbUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return nil, nil, err
	}

	err = MigrateDb(dbUrl)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to migrate the: %v\n", err)
		return nil, nil, err
	}

	return conn, func() error { return conn.Close(ctx) }, nil
}
