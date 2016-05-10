package main

import (
	"github.com/ChimeraCoder/anaconda"
	"os"
	"log"
	"github.com/joho/godotenv"
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"encoding/base64"
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

	mmw := &MyMainWindow{
		api:api,
		modeCtrl:false,
		binaries:make([]string, 0),
	}
	mw := MainWindow{
		AssignTo: &mmw.MainWindow,
		Title: "txt",
		Size   : Size{300, 200},
		Layout: VBox{MarginsZero:true, SpacingZero:true},
		Children: []Widget{
			TextEdit{
				AssignTo: &mmw.tweet,
				OnKeyUp: func(key walk.Key) {
					switch key {
					case walk.KeyControl:
						mmw.modeCtrl = false
					}
				},
				OnKeyPress: func(key walk.Key) {
					switch key {
					case walk.KeyReturn:
						if (mmw.modeCtrl) {
							mmw.post()
							mmw.reset()
						}
					case walk.KeyF1:
						tl := &TimelineWindow{api:api}
						tl.displayTimeline()
					}
				},
				OnKeyDown: func(key walk.Key) {
					switch key {
					case walk.KeyControl:
						mmw.modeCtrl = true
					}
				},
			},
		},
		OnDropFiles: func(path []string) {
			for _, v := range path {
				file, _ := os.Open(v)
				fi, _ := file.Stat()
				size := fi.Size()

				data := make([]byte, size)
				file.Read(data)
				mmw.binaries = append(mmw.binaries, base64.StdEncoding.EncodeToString(data))
				file.Close()
			}
		},
	}

	if _, err := mw.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
