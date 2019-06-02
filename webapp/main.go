package main

import (
	"html/template"
	"net/http"

	"github.com/winstonjay/seizeTheBall/model"

	"google.golang.org/appengine"
)

var indexTemplate = template.Must(template.ParseFiles("index.html"))

func main() {
	http.HandleFunc("/", indexHandler)
	appengine.Main()
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	db, err := model.Connect()
	if err != nil {
		panic(err)
	}
	p, err := model.CurrentBallOwner(db)
	if err != nil {
		panic(err)
	}
	indexTemplate.Execute(w, p)
}
