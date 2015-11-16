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

	queue        []scraper.Article
	openRequests map[string]chan int
}

// NewJefe creates a new jefe with an unstarted scheduler.
func NewJefe() Jefe {
	return Jefe{s: scheduler.MakeScheduler(5, 5), openRequests: make(map[string]chan int)}
}

// Start the scheduler.
func (j *Jefe) Start() {
	j.s.Start()
}

func (j *Jefe) Stop() {
	go j.s.Stop()

	// close all the open requests
	for _, c := range j.openRequests {
		close(c)
	}
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
		j.updateStatus(next.GetLink(), ARTICLE_SENT)
		fmt.Println("next is:", next.GetLink())
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
//
func (j *Jefe) postHandle(rw http.ResponseWriter, request *http.Request) {

	js, err := ioutil.ReadAll(request.Body)
	if err != nil {
		panic(err)
	}

	response := netScraper.Response{}
	err = json.Unmarshal(js, &response)

	fmt.Println("response:", response)

	if response.Data != "" {
		// if got good response
		j.updateStatus(response.Url, ARTICLE_OK)
		close(j.openRequests[response.Url])
		delete(j.openRequests, response.Url)
	} else {
		fmt.Println("response for article")
		// tell the article schedulable that is needs to re-add
		j.updateStatus(response.Url, ARTICLE_BAD)

	}
	rw.WriteHeader(http.StatusOK)
}

// Add an article to the ready queue. The article will be scraped by a
// client then sent back up to the jefe.
func (j *Jefe) AddArticle(article scraper.Article, c chan int) {
	j.queue = append(j.queue, article)
	j.openRequests[article.GetLink()] = c
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

func (j Jefe) updateStatus(name string, status int) {
	if j.openRequests[name] == nil {
		return
	}
	select {

	case j.openRequests[name] <- status:

	default:
		panic("could not send status")
	}
}
