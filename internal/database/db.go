package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
)

func Connect() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("HOST"), os.Getenv("PORT"), os.Getenv("USER"),
		os.Getenv("PASSWORD"), os.Getenv("DBNAME"), os.Getenv("SSLMODE"))

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	return db, db.Ping()
}
