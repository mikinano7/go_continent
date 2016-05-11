package main

import (
	"github.com/lxn/walk"
	"net/url"
)

type Tweet struct {
	ScreenName string
	Text       string
}

type TweetModel struct {
	walk.TableModelBase
	window *TimelineWindow
	items  []*Tweet
}

// ないと死ぬ
func (m *TweetModel) RowCount() int {
	return len(m.items)
}

// ないと死ぬ
func (m *TweetModel) Value(row, col int) interface{} {
	item := m.items[row]

	switch col {
	case 0:
		return item.ScreenName
	case 1:
		return item.Text
	}

	panic("unexpected col")
}

func (m *TweetModel) ResetRows() {
	m.window.tweets, _ = m.window.api.GetHomeTimeline(url.Values{})
	m.items = make([]*Tweet, len(m.window.tweets))

	for i := range m.items {
		m.items[i] = &Tweet{
			ScreenName: m.window.tweets[i].User.ScreenName,
			Text:       m.window.tweets[i].Text,
		}
	}
}
