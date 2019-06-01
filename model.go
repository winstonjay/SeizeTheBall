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
	UserID     int       `json:"userID"`
	TwitterID  string    `json:"twitterID"`
	ScreenName string    `json:"screenName"`
	CreatedAt  time.Time `json:"createdAt"`
}

func (u *User) String() string {
	return fmt.Sprintf("User(UserID=%d TwitterID=%s ScreenName=%s CreatedAt=%s",
		u.UserID, u.TwitterID, u.ScreenName, u.CreatedAt)
}

// Possession : is
type Possession struct {
	PossessionID int       `json:"seizeID"`
	User         User      `json:"user"`
	Start        time.Time `json:"start"`
	End          time.Time `json:"end"`
	Duration     int       `json:"duration"`
}

func ballSeize(db *sql.DB, twitterID, screenName string) error {
	// TODO: test this, and all the rabbit holes it creates
	// havent tried to run at all yet.
	_, err := db.Exec(
		`update user
		set end = now(), duration = timestampdiff(second, start, now())
		where user_id = max(user_id)`)
	if err != nil {
		return err
	}
	return createPossession(db, twitterID, screenName)
}

func createPossession(db *sql.DB, twitterID, screenName string) error {
	userID, err := GetOrCreateUser(db, twitterID, screenName)
	if err != nil {
		return err
	}
	_, err = db.Exec(`insert into user (user_id) values (?)`, userID)
	return err
}

// GetAllPossessions : Get all possessions in database
func GetAllPossessions(db *sql.DB) ([]Possession, error) {
	q := `select (
		p.possession_id, p.Start, p.End, p.Duration,
		u.user_id, u.twitter_id, u.screen_name, u.created_at
	)
	from Possession as p
	inner join user as u on p.user_id = u.user_id`
	res, err := db.Query(q)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	var possessions []Possession
	for res.Next() {
		var p Possession
		err = res.Scan(
			&p.PossessionID, &p.Start, &p.End, &p.Duration,
			&p.User.UserID, &p.User.TwitterID, &p.User.ScreenName, &p.User.CreatedAt)
		if err != nil {
			return nil, err
		}
		possessions = append(possessions, p)
	}
	return possessions, nil
}

// GetOrCreateUser : get user from database by TwitterID and return its userID
// or if user with the given TwitterID does not exist create an new user and
// return the newly created users UserID.
func GetOrCreateUser(db *sql.DB, twitterID, screenName string) (int, error) {
	userID, err := GetUserID(db, twitterID)
	if err != nil {
		return 0, err
	}
	// If we got a user and it has an ID set return its UserID
	if userID != 0 {
		return userID, nil
	}
	// If we didnt get a user we have to create one.
	return CreateUser(db, twitterID, screenName)
}

// CreateUser : create a new user in the database and return the inserted UserID
func CreateUser(db *sql.DB, twitterID, screenName string) (int, error) {
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
	return int(lastID), err
}

// GetAllUsers : get all Users from the database ordered by userID
func GetAllUsers(db *sql.DB) ([]User, error) {
	res, err := db.Query(`select * from user order by user_id`)
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

// GetUserID : get UserID from database via TwitterID
func GetUserID(db *sql.DB, twitterID string) (int, error) {
	var userID int
	err := db.QueryRow(
		`select user_id from user where twitter_id = ?`,
		twitterID).Scan(&userID)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	return userID, nil
}

func getUser(db *sql.DB, userID int) (User, error) {
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
