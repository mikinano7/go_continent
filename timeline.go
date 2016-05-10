package main

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"net/url"
	"os"
)

type TimelineWindow struct {
	*walk.MainWindow
	tweetWindow *TweetWindow
	timeline    *walk.ListBox
	api         *anaconda.TwitterApi
	tweets      []anaconda.Tweet
}

func (tlw *TimelineWindow) Display() {
	list := make([]string, 0)
	tlw.tweets, _ = tlw.api.GetHomeTimeline(url.Values{})
	for _, tweet := range tlw.tweets {
		list = append(list, tweet.Text)
	}

	mw := MainWindow{
		AssignTo: &tlw.MainWindow,
		Title:    "tl",
		Size:     Size{500, 200},
		Layout:   VBox{MarginsZero: true, SpacingZero: true},
		Children: []Widget{
			ListBox{
				AssignTo: &tlw.timeline,
				Model:    list,
				OnKeyPress: func(key walk.Key) {
					switch key {
					case walk.KeyReturn:
						tlw.setReplyStatus()
					}
				},
			},
		},
	}
	if _, err := mw.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func (tlw *TimelineWindow) setReplyStatus() {
	idx := tlw.timeline.CurrentIndex()
	tlw.tweetWindow.reset()
	tlw.tweetWindow.params = url.Values{"in_reply_to_status_id": []string{tlw.tweets[idx].IdStr}}
	tlw.tweetWindow.tweet.SetText(fmt.Sprintf("@%s ", tlw.tweets[idx].User.ScreenName))
}
