package server

import (
	"github.com/opinionated/scheduler"
	"github.com/opinionated/scraper-core/net"
	"github.com/opinionated/scraper-core/scraper"
	"github.com/opinionated/utils/log"
	"sync"
)

// Jefe manages scraping articles. The server should add any RSS feeds that need
// scraping to the Jefe when it starts. SchedulableRSS feeds schedule article
// tasks with the Jefe's scheduler. SchedulableArticle 'manages' the article,
// from adding it to the jefe's work queue to the article being successfully
// scraped. The schedulableArticle puts the article back into the work queue if
// the article does not come back in a reasonable ammount of time or comes back
// with an error. Once the article has been successfully scraped, the Jefe
// passes it off to be analyzed.
// The server treats the Jefe as a queue of articles, and the Jefe is unaware of
// the server implementation.
type Jefe struct {
	// manage RSS pings, manage adds
	s *scheduler.Scheduler

	queue []scraper.Article // articles ready to scrape

	// articles sent and not yet received back
	// article schedulables wait for a response on the chan
	// if the article is not scraped in a reasonable amount of time then
	// the article schedulable will automatically re-queue the article
	openRequests map[string]chan int
	// holds open articles
	openArticles map[string]scraper.Article

	// control read/write to the queue
	mutex *sync.Mutex
}

// NewJefe creates a new Jefe with an unstarted scheduler.
func NewJefe() Jefe {
	return Jefe{s: scheduler.MakeScheduler(5, 5),
		openRequests: make(map[string]chan int),
		openArticles: make(map[string]scraper.Article),
		mutex:        &sync.Mutex{}}
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
	for url, c := range j.openRequests {
		close(c)
		delete(j.openRequests, url)
		delete(j.openArticles, url)
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

	log.Info("handling response:", response)

	_, isOpen := j.openRequests[response.URL]
	if isOpen && response.Error == netScraper.ResponseOk {

		// pass this article off
		// TODO: put in some error checking on the article body
		article := j.openArticles[response.URL]
		article.SetData(response.Data)
		go handleScrapedArticle(article)

		// got a good response
		j.updateStatus(response.URL, ARTICLE_OK)

		// close everything up
		close(j.openRequests[response.URL])
		delete(j.openRequests, response.URL)
		delete(j.openArticles, response.URL)

	} else if isOpen {
		// tell the article schedulable that is needs to re-add
		// don't remove it from the openRequests in case the article
		// comes back before it gets added again
		j.updateStatus(response.URL, ARTICLE_BAD)
	}
	// else a response has already come back, fall through

	return nil
}

// NextArticle returns the next article to scrape. If there is an article to
// scrape, it returns the article and true, else it returns nil and false. Tells
// the schedulableArticle that the article has been sent
func (j *Jefe) NextArticle() (scraper.Article, bool) {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	// look for the next article that has an open requests (the jefe may recieve
	// an article after it has been requeued)
	for j.hasNext() {

		next := j.pop()
		if _, ok := j.openRequests[next.GetLink()]; ok {

			j.updateStatus(next.GetLink(), ARTICLE_SENT)
			log.Info("going to send article", next.GetLink())

			return next, true
		}
	}

	return nil, false
}

// AddArticle adds an article to the ready queue. The article will be scraped by
// a client then sent back up to the Jefe. The chan signals back to the
// schedulable article
func (j *Jefe) AddArticle(article scraper.Article, c chan int) {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	log.Info("adding article", article.GetLink(), "to ready queue")
	j.queue = append(j.queue, article)

	j.openRequests[article.GetLink()] = c
	j.openArticles[article.GetLink()] = article
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
		log.Error("could not send status for article", name)
	}
}
