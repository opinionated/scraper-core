package fetcher_test

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"io/ioutil"
	"net/http"
	"scraper/fetcher"
	mock_rss "scraper/mocks/mock_rss"
	"scraper/scheduler"
	"testing"
	"time"
)

func TestSchedulableRSS(t *testing.T) {
	t.Skip("skipping RSS test for now")
	s := scheduler.MakeScheduler(5, 3)
	s.Start()

	rss := fetcher.CreateSchedulableRSS(&fetcher.WSJRSS{}, 0)
	s.AddSchedulable(rss)
	time.Sleep(time.Duration(6) * time.Second)
	s.Stop()
}

func RunRSS(rss fetcher.RSS) {
	resp, err := http.Get(rss.GetLink())
	fmt.Println("link is:", rss.GetLink())
	if err != nil {
		// TODO: error checking here
		fmt.Println("error getting RSS:", err)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error getting RSS:", err)
		// TODO: error handling
		return
	}

	err = fetcher.GetStories(rss, body)
	if err != nil {
		fmt.Println("error getting RSS:", err)
		return
	}
}

func SignalDone(s *scheduler.Scheduler, c chan bool) {
	<-c
	time.Sleep(time.Duration(2) * time.Second)
	s.Stop()
}

func TestSchedulableRSSMock(t *testing.T) {
	// simulates actual run behavior using mocks
	// used gomock for the mocking
	// generated mock with command:
	// mockgen scraper/fetcher RSS,RSSChannel | tee src/scraper/mocks/mock_rss/mock_rss.go

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// build an actual RSS and push its signals onto the mock
	// use the mock to control what stories get run until code is written to manage it
	wsj := &fetcher.WSJRSS{}
	RunRSS(wsj)

	// build mock RSS, limit how many times it can run
	rss := mock_rss.NewMockRSS(ctrl)
	rss.EXPECT().GetLink().Return(wsj.GetLink()).Times(2)
	rss.EXPECT().GetLink().Return("").AnyTimes()

	// set up and integrate mock chan
	rssChan := mock_rss.NewMockRSSChannel(ctrl)
	rss.EXPECT().GetChannel().Return(rssChan).AnyTimes()
	// don't skip printing stories
	// TODO: remove print found stories
	rssChan.EXPECT().GetNumArticles().Return(0).AnyTimes()

	c := make(chan bool) // signals SignalDone to kill scheduler

	// feed it articles from the real RSS to make it seem like the fake is getting
	// articles each time
	// TODO: look into making the helper functions for getting articles/stories part
	// of the interface
	rssChan.EXPECT().GetArticle(gomock.Any()).Return(wsj.GetChannel().GetArticle(0))
	rssChan.EXPECT().GetArticle(gomock.Any()).Do(func(v interface{}) {
		c <- true
	}).Return(wsj.GetChannel().GetArticle(1))

	// create scheduler and run real schedulable with rss mock
	s := scheduler.MakeScheduler(5, 5)
	go s.AddSchedulable(fetcher.CreateSchedulableRSS(rss, 4))

	// send signal done to loop in the background until we are done with the test
	go SignalDone(s, c)

	// run scheduler in this thread to keep the test cleaner
	s.Run()
}
