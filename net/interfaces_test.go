package netScraper_test

import (
	"fmt"
	"github.com/opinionated/scraper-core/net"
	"github.com/opinionated/scraper-core/net/client"
	"github.com/opinionated/scraper-core/net/server"
	"github.com/opinionated/scraper-core/scraper"
	"net/http"
	"testing"
	"time"
)

//
// NOTE: run tests one at a time until we have graceful server shutdown
//
//

func StartServer(j *server.Jefe) {
	http.HandleFunc("/", j.Handle())
	http.ListenAndServe(":8080", nil)
}

func doSimpleServer() {

	j := server.NewJefe()

	// get the jefe going
	go StartServer(&j)

	// make the scheduler loop quick
	j.SetCycleTime(1)
	j.Start()

	// helper function to add tasks to jefe
	b := func(name string, delay int) {
		article := scraper.WSJArticle{Link: name}
		t := server.CreateSchedulableArticle(&article, delay, &j)
		j.AddSchedulable(t)
	}

	// add three tasks
	// this is the behavior an RSS might give
	b("one", 0)
	b("two", 3)
	b("three", 6)

}

func PanicErr(e error) {
	if e != nil {
		panic(e)
	}
}

func TestServerSimple(t *testing.T) {
	// acts like a client
	// starts the server then adds 3 articles
	// the client then reads the articles faster than
	// they come out, so it looks like there is no work there

	// launch the server
	doSimpleServer()

	expectedUrls := []string{"one", "two", "three"}
	for _, expected := range expectedUrls {

		request, err := client.Get()
		PanicErr(err)
		if request.URL != expected {
			t.Errorf("expected:", expected, "recieved:", request.URL, "\n")
		}

		empty, err := client.Get()
		PanicErr(err)
		if !netScraper.IsEmptyRequest(empty) {
			t.Errorf("expected empty request, got:", empty.URL)
		}

		if request.URL != "" {
			response := netScraper.Response{URL: request.URL, Data: "data", Error: netScraper.ResponseOk}
			client.Post(response)
		}

		time.Sleep(time.Duration(4) * time.Second)
	}
}

func TestServerBad(t *testing.T) {
	// drop a connection
	// launch the server
	doSimpleServer()

	seenTwo := false
	expectedUrls := []string{"one", "two", "three", "two"}
	for _, expected := range expectedUrls {

		request, err := client.Get()
		PanicErr(err)
		if request.URL != expected {
			t.Errorf("expected:", expected, "recieved:", request.URL, "\n")
		}

		fmt.Println("got request:", request.URL)
		if request.URL != "" {
			response := netScraper.Response{URL: request.URL, Data: "data", Error: netScraper.ResponseOk}
			if request.URL == "two" && !seenTwo {
				seenTwo = true
				response.Error = netScraper.ResponseBad
			}
			client.Post(response)
		}

		time.Sleep(time.Duration(4) * time.Second)
	}
}

func TestServerDrop(t *testing.T) {
	// make a bad response
	// launch the server
	doSimpleServer()

	seenTwo := false
	expectedUrls := []string{"one", "two", "three", "two"}
	for _, expected := range expectedUrls {

		request, err := client.Get()
		PanicErr(err)
		if request.URL != expected {
			t.Errorf("expected:", expected, "recieved:", request.URL, "\n")
		}

		fmt.Println("got request:", request.URL)
		if request.URL != "" {
			response := netScraper.Response{URL: request.URL, Data: "data", Error: netScraper.ResponseOk}
			if request.URL == "two" && !seenTwo {
				seenTwo = true
				// don't send
			} else {
				_ = client.Post(response)
			}
		}

		time.Sleep(time.Duration(4) * time.Second)
	}
}
