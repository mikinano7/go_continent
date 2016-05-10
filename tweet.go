package main

import (
	"encoding/base64"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"net/url"
	"os"
	"strings"
	"sync"
)

type TweetWindow struct {
	*walk.MainWindow
	tweet    *walk.TextEdit
	api      *anaconda.TwitterApi
	modeCtrl bool
	binaries []string
}

func (tw *TweetWindow) display() {
	tw.modeCtrl = false
	tw.binaries = make([]string, 0)

	mw := MainWindow{
		AssignTo: &tw.MainWindow,
		Title:    "tw",
		Size:     Size{300, 200},
		Layout:   VBox{MarginsZero: true, SpacingZero: true},
		Children: []Widget{
			TextEdit{
				AssignTo: &tw.tweet,
				OnKeyUp: func(key walk.Key) {
					switch key {
					case walk.KeyControl:
						tw.modeCtrl = false
					}
				},
				OnKeyPress: func(key walk.Key) {
					switch key {
					case walk.KeyReturn:
						if tw.modeCtrl {
							tw.post()
							tw.reset()
						}
					case walk.KeyF1:
						tl := &TimelineWindow{api: tw.api}
						tl.display()
					}
				},
				OnKeyDown: func(key walk.Key) {
					switch key {
					case walk.KeyControl:
						tw.modeCtrl = true
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
				tw.binaries = append(tw.binaries, base64.StdEncoding.EncodeToString(data))
				file.Close()
			}
		},
	}

	if _, err := mw.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func (mmw *TweetWindow) post() {
	var mediaIds []string
	requestParams := url.Values{}
	receiver := mmw.uploadMedia(mmw.binaries)
	for {
		receive, done := <-receiver
		if !done {
			if len(mediaIds) > 0 {
				requestParams.Add("media_ids", strings.Join(mediaIds, ","))
			}
			_, err := mmw.api.PostTweet(mmw.tweet.Text(), requestParams)
			if err != nil {
				fmt.Println(err.Error())
			}
			return
		}
		mediaIds = append(mediaIds, receive)
	}
}

func (mmw *TweetWindow) uploadMedia(binaries []string) <-chan string {
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

func (mmw *TweetWindow) reset() {
	mmw.tweet.SetText("")
	mmw.binaries = make([]string, 0)
}
