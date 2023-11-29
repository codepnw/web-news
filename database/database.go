package database

import (
	"database/sql"
	"os"

	"github.com/upper/db/v4"
)

func ConnectDB() (*sql.DB, error) {
	db, err := sql.Open(os.Getenv("DB_DRIVER"), os.Getenv("DB_DNS"))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func MigrateDB(db db.Session) error {
	script, err := os.ReadFile("./database/migrations/table.sql")
	if err != nil {
		return err
	}

	_, err = db.SQL().Exec(string(script))

	return err
}
