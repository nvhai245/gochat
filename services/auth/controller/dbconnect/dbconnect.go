package dbconnect

import (
	"database/sql"
	"fmt"
	// PostgreSQL driver
	_ "github.com/lib/pq"
)
// ConnectDb PostgreSQL
func ConnectDb() *sql.DB {
	connStr := "postgres://postgres:Harin245@127.0.0.1:5432/gochat?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	// defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("DB connected!")
	return db
}