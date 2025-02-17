package config

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func ConnDB() (*sql.DB, error) {
	sqldb, err := sql.Open(os.Getenv("DB_DRIVER"), os.Getenv("DB_DSN"))
	if err != nil {
		log.Fatal("Cannot connect to db: ", err)
		return nil, err
	}

	return sqldb, nil
}
