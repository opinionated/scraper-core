package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/opinionated/scraper-core/net"
	"io/ioutil"
	"net/http"
	"time"
)

// TODO: handle connection errors
// TODO: scraping logic

type Client struct {
}

// Run runs the client in an infinite loop
func (c *Client) Run() {
	ticker := time.NewTicker(time.Duration(10) * time.Second)

	for {

		// wait for the next itr
		<-ticker.C

		// go get the article
		req, err := Get()
		if err != nil {
			panic(err)
		}

		// don't reply to empty requests
		if netScraper.IsEmptyRequest(req) {
			fmt.Println("got empty request")
			continue
		}

		fmt.Println("client got article with url:", req.URL)

		result := netScraper.Response{URL: req.URL, Data: "", Error: netScraper.ResponseOk}

		err = Post(result)
		if err != nil {
			panic(err)
		}

	}
}

// Get the server for an article to scrape.
func Get() (netScraper.Request, error) {
	c := &http.Client{}

	// get next work unit
	resp, err := c.Get("http://localhost:8080/")
	defer resp.Body.Close()
	if err != nil {
		return netScraper.EmptyRequest(), err
	}

	if resp.StatusCode != 200 {
		fmt.Println("oh nose, did not get OK status:", resp.StatusCode)
		panic("bad status code")
	}

	js, err := ioutil.ReadAll(resp.Body)
	toDo := netScraper.Request{}

	err = json.Unmarshal(js, &toDo)
	if err != nil {
		return netScraper.EmptyRequest(), err

	}
	return toDo, err
}

// Post posts a completed work item up to the server.
func Post(done netScraper.Response) error {

	c := &http.Client{}

	fetchedJSON, err := json.Marshal(done)

	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST",
		"http://localhost:8080/",
		bytes.NewReader(fetchedJSON))

	if err != nil {
		return err
	}

	req.Header.Set("content-type", "application/json")

	postResp, err := c.Do(req)
	defer postResp.Body.Close()
	if err != nil {
		return err
	}

	return nil
}
