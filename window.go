package main

import (
	"github.com/lxn/walk"
	"github.com/ChimeraCoder/anaconda"
	"net/url"
	"strings"
	"fmt"
	"sync"
)

type MyMainWindow struct {
	*walk.MainWindow
	edit *walk.TextEdit
	api  *anaconda.TwitterApi
	modeCtrl bool
	binaries []string
}

func (mmw *MyMainWindow) post() {
	var mediaIds []string
	requestParams := url.Values{}
	receiver := mmw.uploadMedia(mmw.binaries)
	for {
		receive, done := <-receiver
		if !done {
			if len(mediaIds) > 0 {
				url.Values{}
				requestParams.Add("media_ids", strings.Join(mediaIds, ","))
			}
			_, err := mmw.api.PostTweet(mmw.edit.Text(), requestParams)
			if err != nil {
				fmt.Println(err.Error())
			}
			return
		}
		mediaIds = append(mediaIds, receive)
	}
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

func (mmw *MyMainWindow) reset() {
	mmw.edit.SetText("")
	mmw.binaries = make([]string, 0)
}
