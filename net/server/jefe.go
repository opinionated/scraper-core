package server

import (
	"encoding/json"
	"fmt"
	"github.com/opinionated/scheduler/scheduler"
	"github.com/opinionated/scraper-core/net"
	"github.com/opinionated/scraper-core/scraper"
	"io/ioutil"
	"net/http"
)

// Jefe interfaces a scheduler and a task server. Clients ask the jefe for
// articles to scrape. If an article is ready for scraping, it is sent to
// the client. Once the article is scraped it is sent back over to the jefe.
// The jefe makes sure work gets done by re-queueing an article if it didn't
// get scraped in a reasonable ammount of time.
//
// The jefe also manages the scheduler, which in turn manages rss pings.
type Jefe struct {
	// manage RSS pings, manage adds
	s *scheduler.Scheduler

	// TODO: make work queue
	queue []scraper.Article
}

// NewJefe creates a new jefe with an unstarted scheduler.
func NewJefe() Jefe {
	return Jefe{s: scheduler.MakeScheduler(5, 5)}
}

// Start the scheduler.
func (j *Jefe) Start() {
	j.s.Start()
}

// AddSchedulable puts a new task in the scheduler.
func (j *Jefe) AddSchedulable(s scheduler.Schedulable) {
	j.s.AddSchedulable(s)
}

// SetCycleTime sets the scheduler's cycle time
func (j *Jefe) SetCycleTime(cycle int) {
	j.s.SetCycleTime(cycle)
}

// Handle is passed to go's built in http server. It only allows GET and POST
// requests. GET requests ask for an article to scrape and POST requests
// provide the clients result. Actual processing is handed off to internal functions
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

// handle a get request
func (j *Jefe) getHandle(rw http.ResponseWriter, request *http.Request) {
	var work netScraper.Request
	if j.hasNext() {
		next := j.pop()
		work = netScraper.Request{next.GetLink()}
	} else {
		work = netScraper.EmptyRequest()
	}

	bytes, err := json.Marshal(work)
	if err != nil {
		panic(err)
	}

	rw.Write(bytes)
	rw.WriteHeader(http.StatusOK)
}

// handle a POST request
func (j *Jefe) postHandle(rw http.ResponseWriter, request *http.Request) {

	js, err := ioutil.ReadAll(request.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println("js:", string(js))

	fmt.Println("writing post")
	rw.WriteHeader(http.StatusOK)
}

// Add an article to the ready queue. The article will be scraped by a
// client then sent back up to the jefe.
func (j *Jefe) Add(article scraper.Article) {
	j.queue = append(j.queue, article)
}

// for internal use, removes the next article from the queue
func (j *Jefe) pop() scraper.Article {
	ret := j.queue[0]
	j.queue = j.queue[1:]
	return ret
}

// checks if there is an article ready on the queue
func (j Jefe) hasNext() bool {
	return len(j.queue) > 0
}
