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

func InitDb() PgxIface {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	err = MigrateDb(os.Getenv("DATABASE_URL"))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to migrate the: %v\n", err)
		os.Exit(1)
	}

	return conn
}

func CloseDb(ctx context.Context, db PgxIface) func() error {
	return func() error {
		err := db.Close(ctx)
		return err
	}
}
