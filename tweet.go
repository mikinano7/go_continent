package main

import (
	"encoding/base64"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
	"net/url"
	"os"
	"strings"
	"sync"
	"strconv"
	"unicode/utf8"
)

type TweetWindow struct {
	*walk.MainWindow
	tweet    *walk.TextEdit
	count    *walk.ListBox
	api      *anaconda.TwitterApi
	modeCtrl bool
	binaries []string
	params   url.Values
}

func (tw *TweetWindow) Display() {
	tw.modeCtrl = false
	tw.binaries = make([]string, 0)
	tw.params = url.Values{}

	mw := MainWindow{
		AssignTo: &tw.MainWindow,
		Size:     Size{150, 50},
		Layout: VBox{
			Margins:     Margins{0, 0, 0, 0},
			Spacing:     0,
			MarginsZero: true,
			SpacingZero: true,
		},
		Children: []Widget{
			TextEdit{
				AssignTo: &tw.tweet,
				MaxSize:  Size{300, 100},
				OnKeyUp: func(key walk.Key) {
					switch key {
					case walk.KeyControl:
						tw.modeCtrl = false
					default:
						wordCount := utf8.RuneCountInString(tw.tweet.Text())
						tw.count.SetModel([]string{strconv.Itoa(140 - wordCount)})
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
						tl := &TimelineWindow{api: tw.api, tweetWindow: tw}
						tl.Display()
					}
				},
				OnKeyDown: func(key walk.Key) {
					switch key {
					case walk.KeyControl:
						tw.modeCtrl = true
					}
				},
			},
			ListBox{
				AssignTo:&tw.count,
				Model:[]string{strconv.Itoa(140)},
				Enabled:false,
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

	mw.Create()
	tw.removeTitleBar()
	(*mw.AssignTo).Run()
}

func (tw *TweetWindow) post() {
	var mediaIds []string
	receiver := tw.uploadMedia(tw.binaries)
	for {
		receive, done := <-receiver
		if !done {
			if len(mediaIds) > 0 {
				tw.params.Add("media_ids", strings.Join(mediaIds, ","))
			}
			_, err := tw.api.PostTweet(tw.tweet.Text(), tw.params)
			if err != nil {
				fmt.Println(err.Error())
			}
			return
		}
		mediaIds = append(mediaIds, receive)
	}
}

func (tw *TweetWindow) uploadMedia(binaries []string) <-chan string {
	var wg sync.WaitGroup
	receiver := make(chan string)
	go func() {
		for _, v := range binaries {
			wg.Add(1)
			go func(v string) {
				media, _ := tw.api.UploadMedia(v)
				receiver <- media.MediaIDString
				wg.Done()
			}(v)
		}
		wg.Wait()
		close(receiver)
	}()
	return receiver
}

func (tw *TweetWindow) reset() {
	tw.tweet.SetText("")
	tw.binaries = make([]string, 0)
	tw.params = url.Values{}
}

func (tw *TweetWindow) removeTitleBar() {
	wb := &tw.MainWindow.WindowBase
	hWnd := wb.Handle()
	win.SetWindowLong(
		hWnd,
		win.GWL_STYLE,
		win.GetWindowLong(hWnd, win.GWL_STYLE)-win.WS_SYSMENU-win.WS_THICKFRAME-win.WS_CAPTION,
	)
	win.SetWindowPos(
		hWnd,
		win.HWND_TOPMOST,
		0, 0, 0, 0,
		win.SWP_NOMOVE|win.SWP_NOSIZE|win.SWP_DRAWFRAME,
	)
}
