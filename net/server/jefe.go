package server

import (
	"fmt"
	"github.com/opinionated/scheduler/scheduler"
	"github.com/opinionated/scraper-core/net"
	"github.com/opinionated/scraper-core/scraper"
	"sync"
)

// the only behavior this should have in these and other cases is that all articles eventualy get
// scraped (possibly multiple times), but are stored exactly once per article

// Jefe interfaces a scheduler and a task server. It sort of manages a
// state machine for each article. Clients ask the Jefe for articles to
// scrape via http GET. If an article is ready for scraping the Jefe sends
// it to the client. The client tries to scrape the article, then sends the
// results back to the Jefe via http POST. The Jefe automatically requeues
// the article if it doesn't get scraped in a reasonable amount oftime.
//
// The Jefe also manages the scheduler, which in turn manages rss pings.
type Jefe struct {
	// manage RSS pings, manage adds
	s *scheduler.Scheduler

	queue []scraper.Article // articles ready to scrape

	// articles sent and not yet received back
	// article schedulables wait for a response on the chan
	// if the article is not scraped in a reasonable amount of time then
	// the article schedulable will automatically re-queue the article
	openRequests map[string]chan int

	// control read/write to the queue
	mutex *sync.Mutex
}

// NewJefe creates a new Jefe with an unstarted scheduler.
func NewJefe() Jefe {
	return Jefe{s: scheduler.MakeScheduler(5, 5), openRequests: make(map[string]chan int), mutex: &sync.Mutex{}}
}

// Start the scheduler.
func (j *Jefe) Start() {
	j.s.Start()
}

// Stop the scheduler, close all open requests.
func (j *Jefe) Stop() {
	go j.s.Stop()

	j.mutex.Lock()
	// close all the open requests
	for _, c := range j.openRequests {
		close(c)
	}

	j.mutex.Unlock()
}

// AddSchedulable adds a new schedulable to the scheduler.
func (j *Jefe) AddSchedulable(s scheduler.Schedulable) {
	// scheduler is all threadsafe so no need to worry here
	go j.s.Add(s) // spin this out so it won't block
}

// SetCycleTime sets the scheduler's cycle time.
func (j *Jefe) SetCycleTime(cycle int) {
	go j.s.SetCycleTime(cycle)
}

// HandleRepsonse handles a response from the server. It updates
// the article and stores it. Returns an error if there was an
// unexpected issue.
func (j *Jefe) HandleResponse(response netScraper.Response) error {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	fmt.Println("response:", response)

	_, isOpen := j.openRequests[response.URL]
	if isOpen && response.Error == netScraper.ResponseOk {
		// got a good response
		j.updateStatus(response.URL, ARTICLE_OK)

		close(j.openRequests[response.URL])
		delete(j.openRequests, response.URL)

	} else if isOpen {
		fmt.Println("response for article")
		// tell the article schedulable that is needs to re-add
		// don't remove it from the openRequests in case the article
		// comes back before it gets added again
		j.updateStatus(response.URL, ARTICLE_BAD)

	} else {
		// got an article that has already been taken care

	}

	return nil
}

// NextArticle returns the next article to scrape. If there
// is an article to scrape, it returns the article and true,
// else it returns nil and false. Indicates that the article
// has been sent when it returns an article (ie article needs
// to be sent right away).
func (j *Jefe) NextArticle() (scraper.Article, bool) {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	for j.hasNext() {
		// TODO: I think this may still need testing
		next := j.pop()
		if _, ok := j.openRequests[next.GetLink()]; ok {
			j.updateStatus(next.GetLink(), ARTICLE_SENT)
			return next, true
		}
	}

	return nil, false
}

// AddArticle adds an article to the ready queue. The article will be scraped by a
// client then sent back up to the Jefe. The chan signals back to the
// schedulable article
func (j *Jefe) AddArticle(article scraper.Article, c chan int) {
	j.mutex.Lock()
	j.queue = append(j.queue, article)
	j.openRequests[article.GetLink()] = c
	j.mutex.Unlock()
}

// removes an article from the ready queue
func (j *Jefe) pop() scraper.Article {
	// note that exposed methods need to protect access
	ret := j.queue[0]
	j.queue = j.queue[1:]

	return ret
}

// checks if there is an article ready on the queue
func (j Jefe) hasNext() bool {
	// note that exposed methods need to protect access
	return len(j.queue) > 0
}

// helper to push article updates
func (j Jefe) updateStatus(name string, status int) {
	// note that exposed methods need to protect access
	if j.openRequests[name] == nil {
		// if not in open requests
		return
	}

	select {
	case j.openRequests[name] <- status:

	default:
		// this is a big issue because the signal chan should not be blocking
		panic("could not send status")
	}
}
