package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	dbConfig = struct {
		password string
		username string
		hostname string
		schema   string
	}{
		getenv("DB_PASSWORD"),
		getenv("BD_USERNAME"),
		getenv("DB_HOSTNAME"),
		getenv("DB_SCHEMA"),
	}
	connectStr = fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true",
		dbConfig.username, dbConfig.password, dbConfig.hostname, dbConfig.schema)
)

// User : is
type User struct {
	UserID     int64     `json:"userID"`
	TwitterID  string    `json:"twitterID"`
	ScreenName string    `json:"screenName"`
	CreatedAt  time.Time `json:"createdAt"`
}

// Seize : is
type Seize struct {
	SeizeID  int64     `json:"seizeID"`
	User     User      `json:"user"`
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
	Duration int64     `json:"duration"`
}

func seizeBall(db *sql.DB, twitterID, screenName string) error {
	// TODO: test this, and all the rabbit holes it creates
	// havent tried to run at all yet.
	_, err := db.Exec(
		`update user
		set end = now(), duration = timestampdiff(second, start, now())
		where user_id = max(user_id)`)
	if err != nil {
		return err
	}
	return createSeize(db, twitterID, screenName)
}

func createSeize(db *sql.DB, twitterID, screenName string) error {
	userID, err := getOrCreateUserID(db, twitterID, screenName)
	if err != nil {
		return err
	}
	// stmt, err := db.Prepare(`insert into user (user_id) values (?)`)
	// if err != nil {
	// 	return err
	// }
	// defer stmt.Close()
	_, err = db.Exec(`insert into user (user_id) values (?)`, userID)
	return err
}

// MAX(id)

func getOrCreateUserID(db *sql.DB, twitterID, screenName string) (int64, error) {
	userID, err := getUserID(db, twitterID)
	if err != nil {
		return 0, err
	}
	// If we got a user and it has an ID set return its UserID
	if userID != 0 {
		return userID, nil
	}
	// If we didnt get a user we have to create one.
	return createUser(db, twitterID, screenName)
}

// create user and return the inserted ID
func createUser(db *sql.DB, twitterID, screenName string) (int64, error) {
	stmt, err := db.Prepare(
		`insert into user (twitter_id, screen_name) values (?, ?)`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(twitterID, screenName)
	if err != nil {
		return 0, err
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return lastID, err
}

func getAllUsers(db *sql.DB) ([]User, error) {
	res, err := db.Query(`select * from user`)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	users := []User{}
	for res.Next() {
		var u User
		err = res.Scan(&u.UserID, &u.TwitterID, &u.ScreenName, &u.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func getUserID(db *sql.DB, twitterID string) (int64, error) {
	var userID int64
	err := db.QueryRow(
		`select user_id from user where twitter_id = ?`,
		twitterID).Scan(&userID)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	return userID, nil
}

func getUser(db *sql.DB, userID int64) (User, error) {
	var u User
	res, err := db.Query(`select * from user where user_id = ?`, userID)
	if err != nil {
		return u, err
	}
	defer res.Close()
	for res.Next() {
		var u User
		err = res.Scan(&u.UserID, &u.TwitterID, &u.ScreenName, &u.CreatedAt)
		if err != nil {
			return u, err
		}
	}
	return u, nil
}
