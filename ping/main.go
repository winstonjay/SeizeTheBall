package main

import (
	"fmt"

	"github.com/winstonjay/seizeTheBall/model"
)

func main() {
	db, err := model.Connect()
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
