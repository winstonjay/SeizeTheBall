package main

import (
	"html/template"
	"net/http"

	"github.com/winstonjay/seizeTheBall/logger"
	"github.com/winstonjay/seizeTheBall/model"

	"google.golang.org/appengine"
)

var indexTemplate = template.Must(template.ParseFiles("index.html"))

var log = logger.NewLogger()

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
		log.Errorf("Could not connect to the database: %s", err)
	}
	var p model.Possession
	p, err = model.CurrentPossession(db)
	if err != nil {
		log.Errorf("Could not get current ball owner: %s", err)
		p.User.ScreenName = "???"
	}
	indexTemplate.Execute(w, p)
}
