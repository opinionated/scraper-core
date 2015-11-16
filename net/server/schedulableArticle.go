package server

import (
	"fmt"
	"github.com/opinionated/scheduler/scheduler"
	"github.com/opinionated/scraper-core/scraper"
	"math"
	"time"
)

const (
	ARTICLE_SENT = 1
	ARTICLE_OK   = 2
	ARTICLE_BAD  = 3
)

func toString(c int) string {
	if c == ARTICLE_SENT {
		return "SENT"
	}
	if c == ARTICLE_OK {
		return "OK"
	}
	return "BAD"
}

// make an Article schedulable
type SchedulableArticle struct {
	Article scraper.Article
	delay   int
	start   time.Time
	j       *Jefe
	ran     chan int
}

func (task *SchedulableArticle) DoWork(scheduler *scheduler.Scheduler) {
	select {
	case result := <-task.ran:
		fmt.Println("top result for article:", task.Article.GetLink(), "is:", toString(result))

	default:

	}
	fmt.Println("adding article:", task.Article.GetLink())
	task.j.AddArticle(task.Article, task.ran)

	// wait for the article to go off to a client
	res := <-task.ran
	//fmt.Println("wait result for article:", task.Article.GetLink(), "is:", toString(res))
	if res == ARTICLE_OK {
		fmt.Println("article OK from res")
		return
	}

	select {
	case result := <-task.ran:
		fmt.Println("article:", task.Article.GetLink(), "is:", toString(result))
		if result != ARTICLE_BAD {
			return // finish this
		}
		// else fall through to requeue

	case <-time.After(2 * time.Second):
		// fall through to requeue
	}

	fmt.Println("requeueing:", task.Article.GetLink())
	// re-queue
	task.start = time.Now()
	scheduler.AddSchedulable(task)

}

func (task *SchedulableArticle) GetTimeRemaining() int {
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
func CreateSchedulableArticle(task scraper.Article, delay int, j *Jefe) *SchedulableArticle {
	return &SchedulableArticle{task, delay, time.Now(), j, make(chan int, 2)}
}

// check that we implemented this properly
var _ scheduler.Schedulable = (*SchedulableArticle)(nil)
