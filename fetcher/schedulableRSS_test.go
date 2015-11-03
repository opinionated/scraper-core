package fetcher_test

import (
	"github.com/golang/mock/gomock"
	"github.com/opinionated/scheduler/scheduler"
	"github.com/opinionated/scraper-core/fetcher"
	mock_rss "github.com/opinionated/scraper-core/mock_rss"
	"github.com/opinionated/scraper-core/scraper"
	"testing"
	"time"
)

func TestSchedulableRSS(t *testing.T) {
	// t.Skip("skipping RSS test for now")
	s := scheduler.MakeScheduler(5, 3)
	s.Start()

	rss := fetcher.CreateSchedulableRSS(&scraper.WSJRSS{}, 3)
	s.AddSchedulable(rss)
	time.Sleep(time.Duration(6) * time.Second)
	s.Stop()
}

func RunRSS(rss scraper.RSS) {
	err := scraper.UpdateRSS(rss)
	if err != nil {
		panic(err)
	}
}

func SignalDone(s *scheduler.Scheduler, c chan bool) {
	<-c
	time.Sleep(time.Duration(2) * time.Second)
	s.Stop()
}

func TestSchedulableRSSMock(t *testing.T) {
	t.Skip("don't run full test")
	// simulates actual run behavior using mocks
	// used gomock for the mocking
	// generated mock with command:
	// mockgen scraper/fetcher RSS,RSSChannel | tee src/scraper/mocks/mock_rss/mock_rss.go

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// build an actual RSS and push its signals onto the mock
	// use the mock to control what stories get run until code is written to manage it
	wsj := &scraper.WSJRSS{}
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
