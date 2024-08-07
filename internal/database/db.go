package database

import (
	"database/sql"
	_ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {
	psqlInfo := "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable"
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	return db, db.Ping()
}
