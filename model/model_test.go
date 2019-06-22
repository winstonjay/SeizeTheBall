package model

import (
	"database/sql"
	"testing"
)

var (
	// database
	dbUsername = getenv("DATABASE_USERNAME")
	dbPassword = getenv("DATABASE_PASSWORD")
	dbHostname = getenv("DATABASE_HOSTNAME")
	dbSchema   = getenv("DATABASE_SCHEMA")
)

type testInput struct {
	tweetID    string
	screenName string
	twitterID  string
}

var userTestData = []testInput{
	{tweetID: "852021818290352000", twitterID: "123324435", screenName: "abc"},
	{tweetID: "852021818290352100", twitterID: "543213231", screenName: "xyz"},
	{tweetID: "852021818290352110", twitterID: "123132123", screenName: "231434"},
	{tweetID: "852021818290352129", twitterID: "622857704", screenName: "crucial_tech"},
	{tweetID: "852021818290352130", twitterID: "816653", screenName: "TechCrunch"},
	{tweetID: "852021818290352150", twitterID: "18961853", screenName: "nottora2"},
}

func TestConnection(t *testing.T) {
	db, err := Connect(dbUsername, dbPassword, dbHostname, dbSchema)
	if err != nil {
		t.Fatalf("Connection failed: %s", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		t.Fatalf("Connection failed: %s.", err)
	}
}

func TestRegisterPossession(t *testing.T) {
	// to cleanUp all the records we will be creating.
	db := setupTestDB()
	defer cleanUpTestDB(db)
	for i, v := range userTestData {
		err := RegisterPossession(db, v.tweetID, v.twitterID, v.screenName)
		if err != nil {
			t.Errorf("RegisterPossession failed at test %d\n%s", i, err)
		}
	}
	possessions, err := GetAllPossessions(db)
	if err != nil {
		t.Fatalf("Could not get all possessions.\n%s", err)
	}
	if len(possessions) != len(userTestData) {
		t.Errorf("wrong number of possessions created. want=%d got=%d",
			len(userTestData), len(possessions))
	}
	for i, p := range possessions {
		if i == len(possessions)-1 {
			break
		}
		if p.End == nil {
			t.Errorf("Possession.End == nil at test %d", i)
		}
	}
}

func TestCurrentPossession(t *testing.T) {
	db := setupTestDB()
	defer cleanUpTestDB(db)
	for i, v := range userTestData {
		err := RegisterPossession(db, v.tweetID, v.twitterID, v.screenName)
		if err != nil {
			panic(err)
		}
		p, err := CurrentPossession(db)
		if err != nil {
			t.Errorf("CurrentPossession failed at test %d\n%s", i, err)
		}
		if p.TweetID != v.tweetID {
			t.Errorf("wrong tweetID. want=%s got=%s", v.tweetID, p.TweetID)
		}
	}
}

func TestCreatePossesssion(t *testing.T) {
	db := setupTestDB()
	defer cleanUpTestDB(db)
	for i, v := range userTestData {
		err := CreatePossession(db, v.tweetID, v.twitterID, v.screenName)
		if err != nil {
			t.Errorf("CreatePossession failed at test %d\n%s", i, err)
		}
	}
}

func TestGetAllPossesssions(t *testing.T) {
	db := setupTestDB()
	defer cleanUpTestDB(db)
	for _, v := range userTestData {
		err := CreatePossession(db, v.tweetID, v.twitterID, v.screenName)
		if err != nil {
			panic(err)
		}
	}
	possessions, err := GetAllPossessions(db)
	if err != nil {
		t.Fatalf("Could not get all possessions.\n%s", err)
	}
	if len(possessions) != len(userTestData) {
		t.Errorf("wrong number of possessions created. want=%d got=%d",
			len(userTestData), len(possessions))
	}
	for i, p := range possessions {
		if p.TweetID != userTestData[i].tweetID {
			t.Errorf("wrong Possession.TweetID at position %d. want=%s got=%s",
				i, userTestData[i].tweetID, p.TweetID)
		}
		if p.User.TwitterID != userTestData[i].twitterID {
			t.Errorf("wrong Possession.User.TwitterID at position %d. want=%s got=%s",
				i, userTestData[i].twitterID, p.User.TwitterID)
		}
		if p.User.ScreenName != userTestData[i].screenName {
			t.Errorf("wrong Possession.User.ScreenName at position %d. want=%s got=%s",
				i, userTestData[i].screenName, p.User.ScreenName)
		}
	}
}

func TestEndLastPoessession(t *testing.T) {
	db := setupTestDB()
	defer cleanUpTestDB(db)
	for i, v := range userTestData {
		err := CreatePossession(db, v.tweetID, v.twitterID, v.screenName)
		if err != nil {
			panic(err)
		}
		err = EndLastPossession(db)
		if err != nil {
			t.Errorf("EndLastPossession failed at test %d\n%s", i, err)
		}
	}
	possessions, err := GetAllPossessions(db)
	if err != nil {
		t.Fatalf("Could not get all possessions.\n%s", err)
	}
	if len(possessions) != len(userTestData) {
		t.Errorf("wrong number of possessions created. want=%d got=%d",
			len(userTestData), len(possessions))
	}
	for i, p := range possessions {
		if p.End == nil {
			t.Errorf("Possession.End == nil at test %d", i)
		}
	}
}

func TestCreateUser(t *testing.T) {
	db := setupTestDB()
	defer cleanUpTestDB(db)
	testUserIDs := createTestUsers(t, db, userTestData)
	users, err := GetAllUsers(db)
	if err != nil {
		panic(err)
	}
	if len(users) != len(userTestData) {
		t.Fatalf("wrong number of users created. want=%d got=%d", len(userTestData), len(users))
	}
	for i, user := range users {
		if user.UserID != testUserIDs[i] {
			t.Errorf("wrong User.UserID at position %d. want=%d got=%d",
				i, testUserIDs[i], user.UserID)
		}
		if user.ScreenName != userTestData[i].screenName {
			t.Errorf("wrong User.ScreenName at position %d. want=%s got=%s",
				i, userTestData[i].screenName, user.ScreenName)
		}
		if user.TwitterID != userTestData[i].twitterID {
			t.Errorf("wrong User.TwitterID at position %d. want=%s got=%s",
				i, userTestData[i].twitterID, user.TwitterID)
		}
	}
}

func TestGetUserID(t *testing.T) {
	db := setupTestDB()
	defer cleanUpTestDB(db)
	for i, inID := range createTestUsers(t, db, userTestData) {
		twitterID := userTestData[i].twitterID
		outID, err := GetUserID(db, twitterID)
		if err != nil || outID == 0 {
			t.Fatalf("Error getting user where twitterID=%s", twitterID)
		}
		if inID != outID {
			t.Errorf("Wrong UserID where User.TwitterID=%s. want=%d got=%d",
				twitterID, inID, outID)
		}
	}
}

func TestGetOrCreateUser(t *testing.T) {
	db := setupTestDB()
	defer cleanUpTestDB(db)
	// First we will run a round to test if all our users get created.
	testUserIDs := make([]int, len(userTestData))
	for i, testInput := range userTestData {
		newID, err := GetOrCreateUser(db, testInput.twitterID, testInput.screenName)
		if err != nil {
			t.Fatalf("Could not create User: %s", err)
		}
		testUserIDs[i] = newID
	}
	// Then we will run the same loop again but this time no new
	// users should be created and should equal what we just created.
	for i, testInput := range userTestData {
		outID, err := GetOrCreateUser(db, testInput.twitterID, testInput.screenName)
		if err != nil {
			t.Fatalf("Could not create User: %s", err)
		}
		if outID != testUserIDs[i] {
			t.Errorf("wrong User.UserID at position %d. want=%d got=%d",
				i, testUserIDs[i], outID)
		}
	}
	users, err := GetAllUsers(db)
	if err != nil {
		panic(err)
	}
	if len(users) != len(userTestData) {
		t.Fatalf("wrong number of users created. want=%d got=%d", len(userTestData), len(users))
	}
}

func setupTestDB() *sql.DB {
	db, err := Connect(dbUsername, dbPassword, dbHostname, dbSchema)
	if err != nil {
		panic(err)
	}
	return db
}

func cleanUpTestDB(db *sql.DB) {
	_, err := db.Exec(`delete from user`)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`delete from possession`)
	if err != nil {
		panic(err)
	}
	db.Close()
}

func createTestUsers(t *testing.T, db *sql.DB, data []testInput) []int {
	userIDs := make([]int, len(data))
	for i, v := range data {
		lastID, err := CreateUser(db, v.twitterID, v.screenName)
		if err != nil {
			t.Fatalf("Could not create User: %s", err)
		}
		userIDs[i] = lastID
	}
	return userIDs
}
