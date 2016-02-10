package server

import (
	"github.com/opinionated/scheduler"
	"github.com/opinionated/scraper-core/scraper"
	"github.com/opinionated/utils/log"
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

// SchedulableArticle adds an article to the Jefe's ready queue then
// manages the article until it is scraped properly. It will requeue
// the article if needed.
type SchedulableArticle struct {
	article scraper.Article
	delay   int
	start   time.Time
	j       *Jefe
	ran     chan int
}

func (task *SchedulableArticle) Run(scheduler *scheduler.Scheduler) {
	// check if the task ran while we were waiting
	select {
	case result := <-task.ran:
		if result == ARTICLE_OK {
			log.Info("scraped article:", task.article.GetLink())
			return
		}
		if result == ARTICLE_BAD {
			log.Warn("bad result for article:", task.article.GetLink(), "requeueing")
		}
	default:
		// nothing read
	}

	task.j.AddArticle(task.article, task.ran)

	// wait for the article to go off to a client
	res := <-task.ran
	if res == ARTICLE_OK {
		log.Info("scraped article", task.article.GetLink())
		return
	}
	if res == ARTICLE_BAD {
		log.Warn("bad result for article:", task.article.GetLink(), "requeueing")
		// re-queue
		task.start = time.Now()
		task.delay = 15 // set delay to 2 here b/c prev delay was relative
		scheduler.Add(task)
		return
	}

	// once the article is at the client, wait a reasonable amount of time
	// if the article did not come back in the expected time, requeue it
	var waitTime time.Duration = 15

	select {
	case result := <-task.ran:
		if result == ARTICLE_OK {
			log.Info("scraped article", task.article.GetLink())
			return // finish this
		}

		log.Warn("got result", toString(result), "for article", task.article.GetLink(), "requeueing")
		// else fall through to requeue

	case <-time.After(waitTime * time.Second):
		// fall through to requeue
		log.Info("timing out for article", task.article.GetLink())
	}

	task.start = time.Now()
	task.delay = 15 // set delay to 2 here b/c prev delay was relative
	scheduler.Add(task)
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
func CreateSchedulableArticle(task scraper.Article, delay int, j *Jefe) *SchedulableArticle {
	return &SchedulableArticle{task, delay, time.Now(), j, make(chan int, 2)}
}

// check that we implemented this properly
var _ scheduler.Schedulable = (*SchedulableArticle)(nil)
