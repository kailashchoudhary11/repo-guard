package db

import (
	"database/sql"
	"embed"

	"github.com/kailashchoudhary11/repo-guard/initializers"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func Migrate() error {
	initializers.LoadDotEnv()
	database, err := sql.Open("postgres", initializers.LoadDatabaseURl())

	if err != nil {
		return err
	}

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(database, "migrations"); err != nil {
		return err
	}

	return nil
}
