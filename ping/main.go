package main

import (
	"fmt"
	"os"

	"github.com/winstonjay/seizeTheBall/model"
)

var (
	// database
	dbUsername = getenv("DATABASE_USERNAME")
	dbPassword = getenv("DATABASE_PASSWORD")
	dbHostname = getenv("DATABASE_HOSTNAME")
	dbSchema   = getenv("DATABASE_SCHEMA")
)

func main() {
	db, err := model.Connect(dbUsername, dbPassword, dbHostname, dbSchema)
	if err != nil {
		fmt.Printf("Connection failed: %s\n", err)
		return
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		fmt.Printf("Ping failed: %s\n", err)
		return
	}
	fmt.Println("connected successfully")
}

func getenv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		panic("missing required environment variable " + name)
	}
	return v
}
