package main

import (
	"github.com/ChimeraCoder/anaconda"
	"github.com/spf13/viper"
	"log"
	"flag"
)

func main() {
	flag.Parse()
	account := flag.Args()[0]
	viper.SetConfigType("json")
	err := viper.ReadInConfig()
	if err != nil {
		//TODO: なかったらここで作るようにしたい
		log.Fatal("Error loading config file.")
	}

	anaconda.SetConsumerKey(viper.GetString("twitter.consumer.key"))
	anaconda.SetConsumerSecret(viper.GetString("twitter.consumer.secret"))

	//TODO: アカウント設定とか追加して保持したい
	api := anaconda.NewTwitterApi(
		viper.GetString("twitter.credential.token." + account),
		viper.GetString("twitter.credential.secret." + account),
	)

	tw := &TweetWindow{api: api}
	tw.Display()
}
