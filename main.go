package main

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/ChimeraCoder/anaconda"
)

var (
	consumerKey       = getenv("TWITTER_CONSUMER_KEY")
	consumerSecret    = getenv("TWITTER_CONSUMER_SECRET")
	accessToken       = getenv("TWITTER_ACCESS_TOKEN")
	accessTokenSecret = getenv("TWITTER_ACCESS_TOKEN_SECRET")
)

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

func postResponse(api *anaconda.TwitterApi, t anaconda.Tweet) (string, error) {
	newTweetText := fmt.Sprintf("@%s has the ball!", t.User.ScreenName)
	newTweet, err := api.PostTweet(newTweetText, url.Values{})
	if err != nil {
		return fmt.Sprintf("could not tweet '%s': %v", newTweetText, err), err
	}
	return fmt.Sprintf("retweeted %d", newTweet.Id), nil
}

func main() {
	anaconda.SetConsumerKey(consumerKey)
	anaconda.SetConsumerSecret(consumerSecret)
	api := anaconda.NewTwitterApi(accessToken, accessTokenSecret)

	log := newLogger()
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
			continue
		}

		if msg, err := postResponse(api, t); err != nil {
			log.Info(msg)
		} else {
			log.Error(msg)
		}
	}
}
