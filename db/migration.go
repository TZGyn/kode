package db

import (
	"database/sql"
	"embed"
	"fmt"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/pressly/goose/v3"

	_ "modernc.org/sqlite"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func InitDB() error {
	configFilePath, err := xdg.ConfigFile(filepath.Join("kode", "kode.sqlite3"))
	if err != nil {
		return err
	}

	fmt.Println(configFilePath)

	db, err := sql.Open("sqlite", configFilePath)
	if err != nil {
		return err
	}
	defer db.Close()

	// Verify connection
	if err := db.Ping(); err != nil {
		return err
	}

	goose.SetBaseFS(embedMigrations)
	goose.SetVerbose(false)
	goose.SetLogger(goose.NopLogger())
	if err := goose.SetDialect("sqlite"); err != nil {
		return err
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return err
	}

	return nil
}
