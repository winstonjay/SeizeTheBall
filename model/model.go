package model

import (
	"database/sql"
	"fmt"
	"os"
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
	return fmt.Sprintf("User(UserID=%d TwitterID=%s ScreenName=%s CreatedAt=%s)",
		u.UserID, u.TwitterID, u.ScreenName, u.CreatedAt)
}

// Possession : is
type Possession struct {
	PossessionID int        `json:"seizeID"`
	TweetID      string     `json:"tweetID"`
	User         User       `json:"user"`
	Start        *time.Time `json:"start"`
	End          *time.Time `json:"end"`
	Duration     int        `json:"duration"`
}

// Connect :
func Connect() (*sql.DB, error) {
	return sql.Open("mysql", connectStr)
}

// RegisterBallSeize : ...
func RegisterBallSeize(tweetID, twitterID, screenName string) error {
	db, err := sql.Open("mysql", connectStr)
	if err != nil {
		return err
	}
	defer db.Close()
	if err := EndLastPossession(db); err != nil {
		return err
	}
	return CreatePossession(db, tweetID, twitterID, screenName)
}

// CurrentBallOwner : ...
func CurrentBallOwner(db *sql.DB) (Possession, error) {
	var p Possession
	res, err := db.Query(`
	select
		p.possession_id, p.tweet_id, p.start, p.end, p.duration,
		u.user_id, u.twitter_id, u.screen_name, u.created_at
	from (select * from possession order by possession_id desc limit 1) as p
	inner join user as u on u.user_id=p.user_id`)
	if err != nil {
		return p, err
	}
	defer res.Close()
	for res.Next() {
		err = res.Scan(
			&p.PossessionID, &p.TweetID, &p.Start, &p.End, &p.Duration,
			&p.User.UserID, &p.User.TwitterID, &p.User.ScreenName, &p.User.CreatedAt)
		if err != nil {
			return p, err
		}
	}
	return p, nil
}

// CreatePossession : ...
func CreatePossession(db *sql.DB, tweetID, twitterID, screenName string) error {
	userID, err := GetOrCreateUser(db, twitterID, screenName)
	if err != nil {
		return err
	}
	stmt, err := db.Prepare(
		`insert into possession (tweet_id, user_id) values (?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(tweetID, userID)
	return err
}

// EndLastPossession : ...
func EndLastPossession(db *sql.DB) error {
	var lastID int
	err := db.QueryRow(`select max(possession_id) from possession`).Scan(&lastID)
	_, err = db.Exec(
		`update possession
		set end = now(), duration = timestampdiff(second, start, now())
		where possession_id=?`, lastID)
	return err
}

// GetAllPossessions : Get all possessions in database
func GetAllPossessions(db *sql.DB) ([]Possession, error) {
	res, err := db.Query(`
	select
		p.possession_id, p.tweet_id, p.start, p.end, p.duration,
		u.user_id, u.twitter_id, u.screen_name, u.created_at
	from possession as p
	inner join user as u on u.user_id=p.user_id
	order by p.possession_id asc`)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	var possessions []Possession
	for res.Next() {
		var p Possession
		err = res.Scan(
			&p.PossessionID, &p.TweetID, &p.Start, &p.End, &p.Duration,
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
	res, err := db.Query(`select * from user order by user_id asc`)
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

func getenv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		panic("missing required environment variable " + name)
	}
	return v
}
