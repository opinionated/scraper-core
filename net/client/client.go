package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/opinionated/scraper-core/net"
	"io/ioutil"
	"net/http"
)

// GetWork asks the server for an article to scrape.
func GetWork() netScraper.Request {
	c := &http.Client{}

	// get next work unit
	resp, err := c.Get("http://localhost:8080/")
	defer resp.Body.Close()
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 200 {
		fmt.Println("oh nose, did not get OK status:", resp.StatusCode)
		panic("bad status code")
	}

	js, err := ioutil.ReadAll(resp.Body)
	toDo := netScraper.Request{}

	err = json.Unmarshal(js, &toDo)
	if err != nil {
		panic(err)

	}
	return toDo
}

// PostDone posts a completed work item up to the server.
func PostDone(done netScraper.Response) {

	c := &http.Client{}

	fetchedJson, err := json.Marshal(done)
	if err != nil {
		panic(err)
	}
	req, err := http.NewRequest("POST",
		"http://localhost:8080/",
		bytes.NewReader(fetchedJson))

	if err != nil {
		panic(err)
	}
	req.Header.Set("content-type", "application/json")

	postResp, err := c.Do(req)
	defer postResp.Body.Close()
	if err != nil {
		panic(err)
	}
}
