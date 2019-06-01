package main

import (
	"database/sql"
	"testing"
)

type testInput struct {
	screenName string
	twitterID  string
}

var userTestData = []testInput{
	{twitterID: "123324435", screenName: "abc"},
	{twitterID: "543213231", screenName: "xyz"},
	{twitterID: "123132123", screenName: "231434"},
}

func TestConnection(t *testing.T) {
	db, err := sql.Open("mysql", connectStr)
	if err != nil {
		t.Fatalf("Connection failed: %s.\nconnectionStr=%s", err, connectStr)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		t.Fatalf("Connection failed: %s.\nconnectionStr=%s", err, connectStr)
	}
}

func TestCreateUser(t *testing.T) {
	db, err := sql.Open("mysql", connectStr)
	if err != nil {
		t.Fatalf("Connection failed: %s\ncoonectStr=%s", err, connectStr)
	}
	defer cleanUp(db)
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
	db, err := sql.Open("mysql", connectStr)
	if err != nil {
		t.Fatalf("Connection failed: %s\ncoonectStr=%s", err, connectStr)
	}
	defer cleanUp(db)
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
	db, err := sql.Open("mysql", connectStr)
	if err != nil {
		t.Fatalf("Connection failed: %s.\nconnectionStr=%s", err, connectStr)
	}
	defer cleanUp(db)
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

// cleanUpDB : helper function to remove all created entries closin
func cleanUp(db *sql.DB) {
	_, err := db.Exec(`delete from user`)
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
