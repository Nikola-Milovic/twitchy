package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
)

type PgxIface interface {
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	Close(context.Context) error
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
}

func InitDb(ctx context.Context, logger *zap.SugaredLogger) (PgxIface, func() error, error) {
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"))

	// fmt.Println(os.Getenv("POSTGRES_USER"))
	// fmt.Println(os.Getenv("POSTGRES_PASSWORD"))
	// fmt.Println(os.Getenv("POSTGRES_HOST"))
	// fmt.Println(os.Getenv("POSTGRES_PORT"))
	// fmt.Println(os.Getenv("POSTGRES_DB"))
	fmt.Println(dbUrl)

	conn, err := pgx.Connect(context.Background(), dbUrl)

	maxAttempts := 20
	for attempts := 1; attempts <= maxAttempts; attempts++ {
		if err == nil {
			break
		} else {
			if conn != nil {
				err = conn.Ping(ctx)
			} else {
				conn, err = pgx.Connect(context.Background(), dbUrl)
			}
		}

		fmt.Println("Retrying connection to database...")
		fmt.Println(err)
		time.Sleep(time.Duration(attempts) * time.Second)
	}

	if err != nil {
		return nil, nil, err
	}

	err = MigrateDb(dbUrl, logger)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to migrate the: %v\n", err)
		return nil, nil, err
	}

	return conn, func() error { return conn.Close(ctx) }, nil
}
