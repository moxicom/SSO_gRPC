package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var storagePath, migratonsPath, migratonsTable string

	flag.StringVar(&storagePath, "storage-path", "", "path to the storage")
	flag.StringVar(&migratonsPath, "mig-path", "", "path to the migrations")
	flag.StringVar(&migratonsTable, "mig-table", "migrations", "name of the table with migrations")
	flag.Parse()

	if storagePath == "" {
		panic("storage-path is required")
	}

	if migratonsPath == "" {
		panic("mig-math is required")
	}


	m, err := migrate.New(
		"file://" + migratonsPath,
		fmt.Sprintf("sqlite3://%s?x-migrations-table=%s", storagePath, migratonsTable),
	)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err , migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")
			return
		}
		panic(err)
	}

	fmt.Println("migrations applied successfully")
}