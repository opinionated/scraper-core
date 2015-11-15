package netScraper_test

import (
	"github.com/opinionated/scraper-core/net"
	"github.com/opinionated/scraper-core/net/client"
	"github.com/opinionated/scraper-core/net/server"
	"github.com/opinionated/scraper-core/scraper"
	"net/http"
	"testing"
	"time"
)

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

func TestServerSimple(t *testing.T) {
	// acts like a client
	// starts the server then adds 3 articles
	// the client then reads the articles faster than
	// they come out, so it looks like there is no work there

	// launch the server
	doSimpleServer()

	expectedUrls := []string{"one", "two", "three"}
	for _, expected := range expectedUrls {

		request := client.GetWork()
		if request.Url != expected {
			t.Errorf("expected:", expected, "recieved:", request.Url, "\n")
		}

		empty := client.GetWork()
		if !netScraper.IsEmptyRequest(empty) {
			t.Errorf("expected empty request, got:", empty.Url)
		}

		if request.Url != "" {
			response := netScraper.Response{Url: request.Url, Data: "data"}
			client.PostDone(response)
		}

		time.Sleep(time.Duration(4) * time.Second)
	}
}
