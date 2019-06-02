package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/winstonjay/seizeTheBall/model"
)

func handler(w http.ResponseWriter, r *http.Request) {
	db, err := model.Connect()
	if err != nil {
		panic(err)
	}
	users, err := model.GetAllUsers(db)
	if err != nil {
		panic(err)
	}
	s, err := json.Marshal(users)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, string(s))
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
