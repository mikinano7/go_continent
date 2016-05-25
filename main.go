package main

import (
	"github.com/ChimeraCoder/anaconda"
	"github.com/spf13/viper"
	"log"
)

func main() {
	viper.SetConfigFile("./config.json")
	viper.SetConfigType("json")
	err := viper.ReadInConfig()
	if err != nil {
		//TODO: なかったらここで作るようにしたい
		log.Fatal("Error loading config file.")
	}

	anaconda.SetConsumerKey(viper.GetString("twitter.consumer.key"))
	anaconda.SetConsumerSecret(viper.GetString("twitter.consumer.secret"))

	api := anaconda.NewTwitterApi(
		viper.GetString("twitter.credential.token.pn"),
		viper.GetString("twitter.credential.secret.pn"),
	)

	tw := &TweetWindow{api: api}
	tw.Display()
}
