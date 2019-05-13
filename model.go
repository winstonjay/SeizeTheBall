package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var (
	dbConfig = struct {
		password string
		username string
		hostname string
		dbname   string
	}{
		getenv("DB_PASSWORD"),
		getenv("BD_USERNAME"),
		getenv("DB_HOSTNAME"),
		getenv("DB_NAME"),
	}
	connectStr = fmt.Sprintf("%s:%s@tcp(%s)/%s",
		dbConfig.username, dbConfig.password, dbConfig.hostname, dbConfig.dbname)
)

// insert user and return the inserted ID
func insertUser(db *sql.DB, screenname, id string) (int64, error) {
	stmt, err := db.Prepare(
		`insert into user (twitter_screenname, twitter_id) values (?, ?)`)
	if err != nil {
		return -1, err
	}
	res, err := stmt.Exec(screenname, id)
	if err != nil {
		return -1, err
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}
	return lastId, err
}
