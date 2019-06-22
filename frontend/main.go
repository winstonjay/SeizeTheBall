package main

import (
	"html/template"
	"net/http"
	"os"

	"github.com/winstonjay/seizeTheBall/logger"
	"github.com/winstonjay/seizeTheBall/model"

	"google.golang.org/appengine"
)

var (
	// database
	dbUsername = getenv("DATABASE_USERNAME")
	dbPassword = getenv("DATABASE_PASSWORD")
	dbHostname = getenv("DATABASE_HOSTNAME")
	dbSchema   = getenv("DATABASE_SCHEMA")
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
	db, err := model.Connect(dbUsername, dbPassword, dbHostname, dbSchema)
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

func getenv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		panic("missing required environment variable " + name)
	}
	return v
}
