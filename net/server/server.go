package server

import (
	"encoding/json"
	"github.com/opinionated/scraper-core/net"
	"github.com/opinionated/utils/log"
	"io/ioutil"
	"net/http"
)

// ScrapeServer handles routing. The Jefe manages an article work queue.
type ScrapeServer struct {
	j Jefe
}

// TODO: think about placement of this, how to couple it?
// should this bit just be in the main?
func NewScrapeServer() *ScrapeServer {
	j := NewJefe()
	return &ScrapeServer{j}
}

func (s *ScrapeServer) GetJefe() *Jefe {
	return &s.j
}

// Handle is passed to go's built in http server. It only allows GET and POST
// requests. GET requests ask for an article to scrape and POST requests
// provide the clients result. Actual processing is handed off to internal functions.
func (s *ScrapeServer) Handle() func(rw http.ResponseWriter, request *http.Request) {
	return func(rw http.ResponseWriter, request *http.Request) {

		method := request.Method
		if method == "GET" {
			s.getHandle(rw, request)

		} else if method == "POST" {
			s.postHandle(rw, request)
		} else {
			log.Warn("oh nose, unexpected HTTP method:", method)
			rw.WriteHeader(405)
		}
	}
}

// handle a get request (send a scrape request to the client)
func (s *ScrapeServer) getHandle(rw http.ResponseWriter, request *http.Request) {
	// build response
	var work netScraper.Request

	next, has := s.j.NextArticle()

	if has {
		work = netScraper.Request{URL: next.GetLink()}
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
func (s *ScrapeServer) postHandle(rw http.ResponseWriter, request *http.Request) {

	js, err := ioutil.ReadAll(request.Body)
	if err != nil {
		// NOTE: this is not a big server side issue
		panic(err)
	}

	response := netScraper.Response{}
	err = json.Unmarshal(js, &response)

	err = s.j.HandleResponse(response)
	if err != nil {
		panic(err)
	}

	// signal that everything came through alright
	// TODO: is there a case where the client should act differently?
	rw.WriteHeader(http.StatusOK)
}
