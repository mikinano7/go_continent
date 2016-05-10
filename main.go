package main

import (
	"github.com/ChimeraCoder/anaconda"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load("go.env")
	if err != nil {
		//TODO: なかったらここで作るようにしたい
		log.Fatal("Error loading .env file")
	}

	anaconda.SetConsumerKey(os.Getenv("TWITTER_CONSUMER_KEY"))
	anaconda.SetConsumerSecret(os.Getenv("TWITTER_CONSUMER_SECRET"))

	api := anaconda.NewTwitterApi(
		os.Getenv("TWITTER_ACCESS_TOKEN"),
		os.Getenv("TWITTER_ACCESS_SECRET"),
	)

	tw := &TweetWindow{api: api}
	tw.Display()
}
