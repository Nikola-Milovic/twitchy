package db

import (
	"fmt"
	"log"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func MigrateDb(dburl string) error {
	m, err := migrate.New(
		"file://db/migrations",
		dburl,
	)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		if strings.Contains(err.Error(), "no change") {
			fmt.Println(err.Error())
		}
	}

	fmt.Println("Migrated DB")
	return err
}
