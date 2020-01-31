package migrate

import (
	"database/sql"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	_ "github.com/lib/pq"
)

// https://godoc.org/github.com/golang-migrate/migrate#example-NewWithDatabaseInstance
func Migrate(db *sql.DB) {

	// Create postgres specific driver instance
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Create new migration instance
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrate/migrations",
		"nuwc",
		driver,
	)
	if err != nil {
		log.Fatal(err)
	}

	// Migrate all the way up ...
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("An error occurred while syncing the database.. %v", err)
	}
}
