package main

import (
	"github.com/lxn/walk"
	"github.com/ChimeraCoder/anaconda"
	. "github.com/lxn/walk/declarative"
	"fmt"
	"os"
	"net/url"
)

type TimelineWindow struct {
	*walk.MainWindow
	timeline *walk.ListBox
	api  *anaconda.TwitterApi
}

//type Timeline struct {
//	icon ImageView
//	screenName string
//	tweet string
//}

func (tlw *TimelineWindow) displayTimeline() {
	list := make([]string, 0)
	tl, _ := tlw.api.GetHomeTimeline(url.Values{})
	for _, tweet := range tl {
		list = append(list, tweet.Text)
	}

	mw := MainWindow{
		AssignTo: &tlw.MainWindow,
		Title: "tl",
		Size   : Size{500, 200},
		Layout: VBox{MarginsZero:true, SpacingZero:true},
		Children: []Widget{
			ListBox{
				AssignTo: &tlw.timeline,
				Model: list,
			},
		},
	}
	if _, err := mw.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
