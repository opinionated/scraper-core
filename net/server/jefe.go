package main

import (
	"encoding/json"
	"fmt"
	"github.com/opinionated/scheduler/scheduler"
	"github.com/opinionated/scraper-core/net"
	"github.com/opinionated/scraper-core/scraper"
	"io/ioutil"
	"net/http"
)

// Jefe is the server code
type Jefe struct {
	// manage RSS pings, manage adds
	s *scheduler.Scheduler

	// TODO: make work queue
	queue []scraper.Article
}

func NewJefe() Jefe {
	return Jefe{}
}

func (j *Jefe) getHandle(rw http.ResponseWriter, request *http.Request) {

	var work netScraper.Request
	if j.hasNext() {
		next := j.Pop()
		work = netScraper.Request{next.GetLink()}
	} else {
		work = netScraper.EmptyRequest()
	}
	//work.Url = "hello"

	bytes, err := json.Marshal(work)
	if err != nil {
		panic(err)
	}

	rw.Write(bytes)
	rw.WriteHeader(http.StatusOK)
}

func (j *Jefe) postHandle(rw http.ResponseWriter, request *http.Request) {

	js, err := ioutil.ReadAll(request.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println("js:", string(js))

	rw.WriteHeader(http.StatusOK)
}

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

// for queue
func (j *Jefe) Add(article scraper.Article) {
	j.queue = append(j.queue, article)
}

func (j *Jefe) Pop() scraper.Article {
	ret := j.queue[0]
	j.queue = j.queue[1:]
	return ret
}

func (j Jefe) hasNext() bool {
	return len(j.queue) > 0
}
