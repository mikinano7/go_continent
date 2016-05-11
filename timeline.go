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
	timeline    *walk.TableView
	api         *anaconda.TwitterApi
	tweets      []anaconda.Tweet
}

func (tlw *TimelineWindow) Display() {
	model := &TweetModel{window: tlw}
	model.ResetRows()

	mw := MainWindow{
		AssignTo: &tlw.MainWindow,
		Size:     Size{500, 200},
		Layout:   VBox{MarginsZero: true},
		Children: []Widget{
			TableView{
				AssignTo:              &tlw.timeline,
				AlternatingRowBGColor: walk.RGB(255, 255, 200),
				ColumnsOrderable:      false,
				MultiSelection:        false,
				LastColumnStretched:   true,
				Font:                  Font{Family: "Helvetica", PointSize: 10},
				Columns: []TableViewColumn{
					{Title: "Name"},
					{Title: "Tweet"},
				},
				Model: model,
				OnKeyPress: func(key walk.Key) {
					switch key {
					case walk.KeyReturn:
						tlw.setReplyStatus()
					case walk.KeyF5:
						model.ResetRows()
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
