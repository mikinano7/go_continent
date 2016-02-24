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
	modeCtrl bool
	binaries []string
	requestParams url.Values
}

func main() {
	err := godotenv.Load("go.env")
	if err != nil {
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
		requestParams:url.Values{},
	}
	mw := MainWindow{
		AssignTo: &mmw.MainWindow,
		Title: "txt",
		Size   : Size{300, 200},
		Layout: VBox{MarginsZero:true, SpacingZero:true},
		Children: []Widget{
			TextEdit{
				AssignTo: &mmw.edit,
				OnKeyUp: func(key walk.Key) {
					switch key {
					case walk.KeyControl:
						mmw.modeCtrl = false
					}
				},
				OnKeyPress: func(key walk.Key) {
					if (mmw.modeCtrl && key == walk.KeyReturn) {
						mmw.onClicked()
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

func (mmw *MyMainWindow) onClicked() {
	var mediaIds []string
	receiver := mmw.uploadMedia(mmw.binaries)
	for {
		receive, done := <-receiver
		if !done {
			if len(mediaIds) > 0 {
				mmw.requestParams.Add("media_ids", strings.Join(mediaIds, ","))
			}
			_, err := mmw.api.PostTweet(mmw.edit.Text(), mmw.requestParams)
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
	mmw.requestParams = url.Values{}
	mmw.binaries = make([]string, 0)
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
