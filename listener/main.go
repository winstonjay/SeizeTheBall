package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/ChimeraCoder/anaconda"
	"github.com/winstonjay/seizeTheBall/logger"
	"github.com/winstonjay/seizeTheBall/model"
)

var (
	// twitter
	consumerKey       = getenv("TWITTER_CONSUMER_KEY")
	consumerSecret    = getenv("TWITTER_CONSUMER_SECRET")
	accessToken       = getenv("TWITTER_ACCESS_TOKEN")
	accessTokenSecret = getenv("TWITTER_ACCESS_TOKEN_SECRET")

	// database
	dbUsername = getenv("DATABASE_USERNAME")
	dbPassword = getenv("DATABASE_PASSWORD")
	dbHostname = getenv("DATABASE_HOSTNAME")
	dbSchema   = getenv("DATABASE_SCHEMA")
)

const (
	keyphrase  = "i have the ball"
	screenName = "@seizetheball"
)

func main() {
	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	api := anaconda.NewTwitterApi(accessToken, accessTokenSecret)

	log := logger.NewLogger()
	api.SetLogger(log)

	stream := api.PublicStreamFilter(url.Values{
		"track": []string{screenName},
	})

	defer stream.Stop()

	for v := range stream.C {
		t, ok := v.(anaconda.Tweet)
		if !ok {
			log.Warningf("received unexpected value of type %T", v)
			continue
		}
		if t.RetweetedStatus != nil {
			continue
		}
		if !matchTweet(t.Text) {
			log.Infof("recived tweet with unmatching text got='%s'", t.Text)
			continue
		}

		db, err := model.Connect(dbUsername, dbPassword, dbHostname, dbSchema)
		defer db.Close()
		if err != nil {
			log.Errorf("Could not connect to DB: %s", err)
		}
		// Register who has now taken possession of the ball in out database.
		err = model.RegisterPossession(db, t.IdStr, t.User.IdStr, t.User.ScreenName)
		if err != nil {
			log.Errorf("Could not register possession: %s", err)
			continue
		}
		log.Infof("Registered new possession. tweet=%s", t.IdStr)
		// Finally tell twitter who now has possession of the ball.
		newTweetText := fmt.Sprintf("@%s has the ball! 🏆⚽️\n%s",
			t.User.ScreenName, uniqueIDString())
		newTweet, err := api.PostTweet(newTweetText, url.Values{})
		if err != nil {
			log.Errorf("could not tweet '%s': %v", newTweetText, err)
			continue
		}
		log.Infof("Tweeted %d", newTweet.Id)
	}
}

// func registerBallPossession()

func getenv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		panic("missing required environment variable " + name)
	}
	return v
}

func matchTweet(tweetText string) bool {
	return strings.Contains(strings.ToLower(tweetText), keyphrase)
}

func uniqueIDString() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
