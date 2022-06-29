package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func Connect() *sql.DB {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/test_24h")

	if err != nil {
		panic(err.Error())
	}

	return db
}
