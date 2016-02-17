package main

import (
	"github.com/ChimeraCoder/anaconda"
	"os"
	"log"
	"github.com/joho/godotenv"
	"net/url"
	"fmt"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"encoding/base64"
	"sync"
	"strings"
)

type MyMainWindow struct {
	*walk.MainWindow
	edit *walk.TextEdit
	api  *anaconda.TwitterApi
}

var (
	binaries []string
	values = url.Values{}
	modeCtrl bool
)

func main() {
	err := godotenv.Load("go.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	anaconda.SetConsumerKey(os.Getenv("TWITTER_CONSUMER_KEY"))
	anaconda.SetConsumerSecret(os.Getenv("TWITTER_CONSUMER_SECRET"))

	api := anaconda.NewTwitterApi(
		os.Getenv("TWITTER_ACCESS_TOKEN2"),
		os.Getenv("TWITTER_ACCESS_SECRET2"),
	)

	mmw := &MyMainWindow{api:api}
	mw := MainWindow{
		AssignTo: &mmw.MainWindow,
		Title: "仕事しろ",
		Size   : Size{300, 200},
		Layout: VBox{MarginsZero:true, SpacingZero:true},
		Children: []Widget{
			TextEdit{
				AssignTo: &mmw.edit,
				OnKeyDown: func(key walk.Key) {
					if (key == walk.KeyControl) {
						modeCtrl = true
					} else if (key == walk.KeyReturn && modeCtrl) {
						mmw.onClicked()
					}
				},
			},
		},
		OnDropFiles: func(path []string) {
			for _, v := range path {
				file, _ := os.Open(v)
				defer file.Close()
				fi, _ := file.Stat()
				size := fi.Size()

				data := make([]byte, size)
				file.Read(data)
				binaries = append(binaries, base64.StdEncoding.EncodeToString(data))
			}
		},
	}

	if _, err := mw.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func (mmw *MyMainWindow) onClicked() {
	var mediaIds []string
	receiver := mmw.uploadMedia(binaries)
	for {
		receive, done := <-receiver
		if !done {
			if len(mediaIds) > 0 {
				values.Add("media_ids", strings.Join(mediaIds, ","))
			}
			_, err := mmw.api.PostTweet(mmw.edit.Text(), values)
			if err != nil {
				fmt.Println(err.Error())
			}
			reset(mmw)
			return
		}
		mediaIds = append(mediaIds, receive)
	}
}

func reset(mmw *MyMainWindow) {
	mmw.edit.SetText("")
	values = url.Values{}
	binaries = make([]string, 0)
	modeCtrl = false
}

func (mmw *MyMainWindow) uploadMedia(binaries []string) <-chan string {
	var wg sync.WaitGroup
	receiver := make(chan string)
	go func() {
		for _, v := range binaries {
			wg.Add(1)
			go func(v string) {
				media, _ := mmw.api.UploadMedia(v)
				receiver <- media.MediaIDString
				wg.Done()
			}(v)
		}
		wg.Wait()
		close(receiver)
	}()
	return receiver
}
