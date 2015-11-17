package server

import (
	"encoding/json"
	"fmt"
	"github.com/opinionated/scheduler/scheduler"
	"github.com/opinionated/scraper-core/net"
	"github.com/opinionated/scraper-core/scraper"
	"io/ioutil"
	"net/http"
	"sync"
)

// Jefe interfaces a scheduler and a task server. Clients ask the Jefe for
// articles to scrape via http GET. If an article is ready for scraping the
// Jefe sends it to the client. The client tries to scrape the article, then
// sends the results back to the Jefe via http POST. The Jefe automatically
// requeues the article if it doesn't get scraped in a reasonable amount of
// time.
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

	// close all the open requests
	for _, c := range j.openRequests {
		close(c)
	}
}

// AddSchedulable adds a new schedulable to the scheduler.
func (j *Jefe) AddSchedulable(s scheduler.Schedulable) {
	j.s.Add(s)
}

// SetCycleTime sets the scheduler's cycle time.
func (j *Jefe) SetCycleTime(cycle int) {
	go j.s.SetCycleTime(cycle)
}

// Handle is passed to go's built in http server. It only allows GET and POST
// requests. GET requests ask for an article to scrape and POST requests
// provide the clients result. Actual processing is handed off to internal functions.
func (j *Jefe) Handle() func(rw http.ResponseWriter, request *http.Request) {
	return func(rw http.ResponseWriter, request *http.Request) {

		method := request.Method
		if method == "GET" {
			j.getHandle(rw, request)

		} else if method == "POST" {
			j.postHandle(rw, request)
		} else {
			fmt.Println("oh nose, unexpected HTTP method:", method)
			rw.WriteHeader(405)
		}
	}
}

// handle a get request (send a scrape request to the client)
func (j *Jefe) getHandle(rw http.ResponseWriter, request *http.Request) {
	// build response
	var work netScraper.Request

	if j.hasNext() {
		next := j.pop()

		// signal that the article is off for scraping
		j.updateStatus(next.GetLink(), ARTICLE_SENT)

		fmt.Println("next is:", next.GetLink())
		work = netScraper.Request{next.GetLink()}

	} else { // !hasNext
		work = netScraper.EmptyRequest()
	}

	// marshal
	bytes, err := json.Marshal(work)
	if err != nil {
		// this is a big server side issue
		panic(err)
	}

	rw.Write(bytes)
	rw.WriteHeader(http.StatusOK)
}

// handle a POST request (receive a scrape response from the client)
func (j *Jefe) postHandle(rw http.ResponseWriter, request *http.Request) {

	js, err := ioutil.ReadAll(request.Body)
	if err != nil {
		// NOTE: this is not a big server side issue
		panic(err)
	}

	response := netScraper.Response{}
	err = json.Unmarshal(js, &response)
	if err != nil {
		panic(err)
	}

	fmt.Println("response:", response)

	if response.Error == netScraper.ResponseOk {
		// got a good response
		j.updateStatus(response.URL, ARTICLE_OK)

		close(j.openRequests[response.URL])
		delete(j.openRequests, response.URL)

	} else {
		fmt.Println("response for article")
		// tell the article schedulable that is needs to re-add
		j.updateStatus(response.URL, ARTICLE_BAD)

	}

	// signal that everything came through alright
	// TODO: is there a case where the client should act differently?
	rw.WriteHeader(http.StatusOK)
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
	j.mutex.Lock()

	ret := j.queue[0]
	j.queue = j.queue[1:]

	j.mutex.Unlock()

	return ret
}

// checks if there is an article ready on the queue
func (j Jefe) hasNext() bool {
	j.mutex.Lock()
	defer j.mutex.Unlock()

	return len(j.queue) > 0
}

// helper to push article updates
func (j Jefe) updateStatus(name string, status int) {
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
