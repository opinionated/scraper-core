package netScraper_test

import (
	"flag"
	"fmt"
	"github.com/opinionated/scraper-core/net"
	"github.com/opinionated/scraper-core/net/client"
	"github.com/opinionated/scraper-core/net/server"
	"github.com/opinionated/scraper-core/scraper"
	"github.com/opinionated/utils/log"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	flag.Parse()
	log.InitStd()
	os.Exit(m.Run())
}

//
// NOTE: run tests one at a time until we have graceful server shutdown
// these tests use lots of hardcoded values, and may not work if you change
// the timings here or somewhere else.
//
//

// start the servr on local:8080
func StartServer(s *server.ScrapeServer) {
	http.HandleFunc("/", s.Handle())
	http.ListenAndServe(":8080", nil)
}

// spin up the server and add three articles staggered
func doSimpleServer() {

	s := server.NewScrapeServer()

	// get the jefe going
	go StartServer(s)

	// make the scheduler loop quick
	s.GetJefe().SetCycleTime(1)
	s.GetJefe().Start()

	// helper function to add tasks to jefe
	b := func(name string, delay int) {
		article := scraper.WSJArticle{Link: name}
		t := server.CreateSchedulableArticle(&article, delay, s.GetJefe())
		s.GetJefe().AddSchedulable(t)
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
	time.Sleep(1 * time.Second)

	expectedUrls := []string{"one", "two", "three"}
	for _, expected := range expectedUrls {

		fmt.Println("client fetching")
		request, err := client.Get("http://localhost:8080")
		PanicErr(err)
		if request.URL != expected {
			t.Errorf("expected:", expected, "recieved:", request.URL, "\n")
		}

		empty, err := client.Get("http://localhost:8080")
		PanicErr(err)
		if !netScraper.IsEmptyRequest(empty) {
			t.Errorf("expected empty request, got:", empty.URL)
		}

		if request.URL != "" {
			response := netScraper.Response{URL: request.URL, Data: "data", Error: netScraper.ResponseOk}
			client.Post("http://localhost:8080", response)
		}

		time.Sleep(time.Duration(4) * time.Second)
	}
}

func TestServerSingleLatePreRerun(t *testing.T) {
	// respond late to a request
	doSimpleServer()

	seenTwo := false
	time.Sleep(1 * time.Second)
	expectedUrls := []string{"one", "two", "three"}
	for _, expected := range expectedUrls {

		request, err := client.Get("http://localhost:8080")
		PanicErr(err)
		if request.URL != expected {
			t.Errorf("expected:", expected, "recieved:", request.URL, "\n")
		}

		fmt.Println("got request:", request.URL)
		if request.URL != "" {
			response := netScraper.Response{URL: request.URL, Data: "data", Error: netScraper.ResponseOk}
			if request.URL == "two" && !seenTwo {
				seenTwo = true
				time.Sleep(time.Duration(3) * time.Second)
			}
			client.Post("http://localhost:8080", response)
		}

		time.Sleep(time.Duration(4) * time.Second)
	}
}

func TestServerSingleLatePostRerun(t *testing.T) {
	// respond late to a request
	doSimpleServer()

	seenTwo := false
	time.Sleep(1 * time.Second)
	expectedUrls := []string{"one", "two", "three"}
	for _, expected := range expectedUrls {

		request, err := client.Get("http://localhost:8080")
		PanicErr(err)
		if request.URL != expected {
			t.Errorf("expected:", expected, "recieved:", request.URL, "\n")
		}

		fmt.Println("got request:", request.URL)
		if request.URL != "" {
			response := netScraper.Response{URL: request.URL, Data: "data", Error: netScraper.ResponseOk}
			if request.URL == "two" && !seenTwo {
				seenTwo = true
				time.Sleep(time.Duration(6) * time.Second)
			}
			client.Post("http://localhost:8080", response)
		}

		time.Sleep(time.Duration(4) * time.Second)
	}
}

func TestServerMultiLatePostRerun(t *testing.T) {
	// respond late to a request
	doSimpleServer()

	seenTwo := false
	expectedUrls := []string{"one", "two", "three", "two"}
	time.Sleep(1 * time.Second)
	i := 0
	for _, expected := range expectedUrls {

		request, err := client.Get("http://localhost:8080")
		PanicErr(err)
		if request.URL != expected {
			t.Errorf("expected: %s recieved: %s\n", request.URL, expected)
		}
		i++
		fmt.Println("got request:", request.URL)
		if request.URL != "" {
			response := netScraper.Response{URL: request.URL, Data: "data", Error: netScraper.ResponseOk}
			if request.URL == "two" && !seenTwo {
				seenTwo = true
				go func() {
					time.Sleep(time.Duration(12) * time.Second)
					client.Post("http://localhost:8080", response)
				}()
			} else {
				client.Post("http://localhost:8080", response)
			}

		}

		time.Sleep(time.Duration(4) * time.Second)
	}

	if i != len(expectedUrls) {
		t.Errorf("did not get through all the urls")
	}
}

func TestServerBad(t *testing.T) {
	// respond with an error
	// launch the server
	doSimpleServer()

	seenTwo := false
	time.Sleep(1 * time.Second)
	expectedUrls := []string{"one", "two", "three", "two"}
	for _, expected := range expectedUrls {

		request, err := client.Get("http://localhost:8080")
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
			client.Post("http://localhost:8080", response)
		}

		time.Sleep(time.Duration(4) * time.Second)
	}
}

func TestServerDrop(t *testing.T) {
	// drop a request
	// launch the server
	doSimpleServer()

	seenTwo := false
	time.Sleep(1 * time.Second)
	expectedUrls := []string{"one", "two", "three", "two"}
	for _, expected := range expectedUrls {

		request, err := client.Get("http://localhost:8080")
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
				_ = client.Post("http://localhost:8080", response)
			}
		}

		time.Sleep(time.Duration(4) * time.Second)
	}
}
