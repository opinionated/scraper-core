package fetcher

import (
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"scraper/scheduler"
	"time"
)

// make an RSS schedulable
type SchedulableRSS struct {
	RSSFeed RSS
	delay   int
	start   time.Time
}

func (rss *SchedulableRSS) DoWork(scheduler *scheduler.Scheduler) {
	fmt.Println("goint to run RSS")
	resp, err := http.Get(rss.RSSFeed.GetLink())
	if err != nil {
		// TODO: error checking here
		fmt.Println("error getting RSS:", err)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error reading body")
		// TODO: error handling
		return
	}

	err = GetStories(rss.RSSFeed, body)
	if err != nil {
		fmt.Println("error getting stories")
		return
	}

	fmt.Println("OK")
	// TODO: use config file to control the timing here
	toSchedule := CreateSchedulableArticle(rss.RSSFeed.GetChannel().GetArticle(0), 1)
	go scheduler.AddSchedulable(toSchedule)
	if rss.IsLoopable() && scheduler.IsRunning() {
		rss.start = time.Now()
		go scheduler.AddSchedulable(rss)
	}
}

func (rss *SchedulableRSS) GetTimeRemaining() int {
	remainingTime := float64(rss.delay) - time.Since(rss.start).Seconds()
	if remainingTime <= 0 {
		return 0
	}
	return int(math.Ceil(remainingTime))
}

func (rss *SchedulableRSS) IsLoopable() bool {
	// TODO: make this true once out of testing
	return true
}

func (rss *SchedulableRSS) SetTimeRemaining(remaining int) {
	rss.delay = remaining
}

// factory to make schedulable rss
func CreateSchedulableRSS(rss RSS, delay int) *SchedulableRSS {
	return &SchedulableRSS{rss, delay, time.Now()}
}

// check that we implemented this properly
var _ scheduler.Schedulable = (*SchedulableRSS)(nil)
