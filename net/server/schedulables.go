package server

import (
	"fmt"
	"github.com/opinionated/scheduler/scheduler"
	"github.com/opinionated/scraper-core/scraper"
	"math"
	"time"
)

// make an Article schedulable
type SchedulableArticle struct {
	Article scraper.Article
	delay   int
	start   time.Time
	j       *Jefe
}

func (task *SchedulableArticle) DoWork(scheduler *scheduler.Scheduler) {
	fmt.Println("adding article:", task.Article.GetLink())
	task.j.Add(task.Article)
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
	return &SchedulableArticle{task, delay, time.Now(), j}
}

// check that we implemented this properly
var _ scheduler.Schedulable = (*SchedulableArticle)(nil)

//
//
//
// make an RSS schedulable
type SchedulableRSS struct {
	rss         scraper.RSS
	delay       int
	start       time.Time
	oldArticles map[string]bool
	j           *Jefe
}

func (task *SchedulableRSS) DoWork(scheduler *scheduler.Scheduler) {
	fmt.Println("goint to run RSS")

	err := scraper.UpdateRSS(task.rss)
	if err != nil {
		fmt.Println("error getting stories")
		return
	}

	// mark all articles as not in list
	for key := range task.oldArticles {
		task.oldArticles[key] = false
	}

	// schedule any new articles
	// an article is new if it wasn't in the last RSS ping
	delay := 10 // TODO: create legitimate task delays
	for i := 0; i < task.rss.GetChannel().GetNumArticles(); i++ {
		article := task.rss.GetChannel().GetArticle(i)
		if _, inOld := task.oldArticles[article.GetLink()]; !inOld {
			toSchedule := CreateSchedulableArticle(article, delay, task.j)
			delay += 10
			go scheduler.AddSchedulable(toSchedule)
		}

		// add or update what we found
		task.oldArticles[article.GetLink()] = true
	}

	// remove any articles not in the set
	for key, inList := range task.oldArticles {
		if !inList {
			delete(task.oldArticles, key)
		}
	}

	// reschedule this task
	if task.IsLoopable() && scheduler.IsRunning() {
		task.start = time.Now()
		task.rss.GetChannel().ClearArticles()
		go scheduler.AddSchedulable(task)
	}
}

func (task *SchedulableRSS) GetTimeRemaining() int {
	remainingTime := float64(task.delay) - time.Since(task.start).Seconds()
	if remainingTime <= 0 {
		return 0
	}
	return int(math.Ceil(remainingTime))
}

func (task *SchedulableRSS) IsLoopable() bool {
	// TODO: make this true once out of testing
	return true
}

func (task *SchedulableRSS) SetTimeRemaining(remaining int) {
	task.delay = remaining
}

// factory to make schedulable task
func CreateSchedulableRSS(task scraper.RSS, delay int, j *Jefe) *SchedulableRSS {
	return &SchedulableRSS{task, delay, time.Now(), make(map[string]bool), j}
}

// check that we implemented this properly
var _ scheduler.Schedulable = (*SchedulableRSS)(nil)
