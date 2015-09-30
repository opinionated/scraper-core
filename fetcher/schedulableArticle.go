package fetcher

import (
	"fmt"
	"math"
	"scraper/scheduler"
	"time"
)

// make an Article schedulable
type SchedulableArticle struct {
	Article Article
	delay   int
	start   time.Time
}

func (article *SchedulableArticle) DoWork(scheduler *scheduler.Scheduler) {
	fmt.Println("goint to get article")
	err := DoGetArticle(article.Article)
	if err != nil {
		fmt.Println("error getting article")
		return
	}
	fmt.Println("article body is:", article.Article.GetData())
}

func (article *SchedulableArticle) GetTimeRemaining() int {
	remainingTime := float64(article.delay) - time.Since(article.start).Seconds()
	if remainingTime <= 0 {
		return 0
	}
	return int(math.Ceil(remainingTime))
}

func (article *SchedulableArticle) IsLoopable() bool {
	// TODO: make this true once out of testing
	return false
}

func (article *SchedulableArticle) SetTimeRemaining(remaining int) {
	article.delay = remaining
}

// factory to make schedulable article
func CreateSchedulableArticle(article Article, delay int) *SchedulableArticle {
	return &SchedulableArticle{article, delay, time.Now()}
}

// check that we implemented this properly
var _ scheduler.Schedulable = (*SchedulableArticle)(nil)
