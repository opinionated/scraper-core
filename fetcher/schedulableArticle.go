package fetcher

import (
	"fmt"
	"github.com/opinionated/scheduler/scheduler"
	"github.com/opinionated/scraper-core/scraper" // this name creates some ambiguity, get a better name...
	"math"
	"time"
)

// make an Article schedulable
type SchedulableArticle struct {
	Article scraper.Article
	delay   int
	start   time.Time
}

func (task *SchedulableArticle) Run(scheduler *scheduler.Scheduler) {
	fmt.Println("goint to get task")
	err := scraper.ScrapeArticle(task.Article)
	if err != nil {
		fmt.Println("error getting task")
		return
	}
	fmt.Println("task body is:", task.Article.GetData())
}

func (task *SchedulableArticle) TimeRemaining() int {
	remainingTime := float64(task.delay) - time.Since(task.start).Seconds()
	if remainingTime <= 0 {
		return 0
	}
	return int(math.Ceil(remainingTime))
}

func (task *SchedulableArticle) IsLoopable() bool {
	// TODO: make this true once out of testing
	return false
}

func (task *SchedulableArticle) SetTimeRemaining(remaining int) {
	task.delay = remaining
}

// factory to make schedulable task
func CreateSchedulableArticle(task scraper.Article, delay int) *SchedulableArticle {
	return &SchedulableArticle{task, delay, time.Now()}
}

// check that we implemented this properly
var _ scheduler.Schedulable = (*SchedulableArticle)(nil)
