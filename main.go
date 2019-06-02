package main

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/ChimeraCoder/anaconda"
	"github.com/sirupsen/logrus"
)

var (
	consumerKey       = getenv("TWITTER_CONSUMER_KEY")
	consumerSecret    = getenv("TWITTER_CONSUMER_SECRET")
	accessToken       = getenv("TWITTER_ACCESS_TOKEN")
	accessTokenSecret = getenv("TWITTER_ACCESS_TOKEN_SECRET")
)

func main() {
	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	api := anaconda.NewTwitterApi(accessToken, accessTokenSecret)

	log := &logger{logrus.New()}
	api.SetLogger(log)

	stream := api.PublicStreamFilter(url.Values{
		"track": []string{"@whohastheball_"},
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

		// Register who has now taken possession of the ball in out database.
		err := RegisterBallSeize(t.IdStr, t.User.IdStr, t.User.ScreenName)
		if err != nil {
			log.Error(err)
			return
		}
		// Finally tell twitter who now has possession of the ball.
		newTweetText := fmt.Sprintf("@%s has the ball! 🏆⚽️", t.User.ScreenName)
		newTweet, err := api.PostTweet(newTweetText, url.Values{})
		if err != nil {
			log.Errorf("could not tweet '%s': %v", newTweetText, err)
			return
		}
		log.Infof("Tweeted %d", newTweet.Id)
	}
}

func getenv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		panic("missing required environment variable " + name)
	}
	return v
}

func matchTweet(tweetText string) bool {
	return strings.Contains(strings.ToLower(tweetText), "i have the ball")
}
