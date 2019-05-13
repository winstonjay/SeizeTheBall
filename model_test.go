package main

import (
	"database/sql"
	"testing"
)

func TestConnection(t *testing.T) {
	db, err := sql.Open("mysql", connectStr)
	if err != nil {
		t.Fatalf("ERROR: %s", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		t.Fatalf("ERROR: %s", err)
	}
}

func TestInsertUser(t *testing.T) {
	db, err := sql.Open("mysql", connectStr)
	if err != nil {
		t.Fatalf("ERROR: %s", err)
	}
	defer db.Close()
	lastId, err := insertUser(db, "hastheball_", "122345432")
	if err != nil {
		t.Fatalf("Could not insert user: %s", err)
	}
	t.Log(lastId)
}
