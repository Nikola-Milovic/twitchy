package db

import (
	"log"
	"os"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"
)

func MigrateDb(dburl string, l *zap.SugaredLogger) error {
	m, err := migrate.New(
		"file:///"+os.Getenv("MIGRATION_PATH"),
		dburl,
	)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		if strings.Contains(err.Error(), "no change") {
			l.Error(err.Error())
		}
	}

	l.Info("Migrated DB")
	return err
}
