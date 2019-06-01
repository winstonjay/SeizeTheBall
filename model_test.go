package main

import (
	"database/sql"
	"testing"
)

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
	_, err = createUser(db, "hastheball_", "122345432")
	if err != nil {
		t.Fatalf("Could not insert user: %s", err)
	}
	users, err := getAllUsers(db)
	if err != nil {
		panic(err)
	}
	if len(users) != 1 {
		t.Fatalf("User not created: %d exist", len(users))
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
