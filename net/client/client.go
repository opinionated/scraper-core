package client

import (
	"bytes"
	"encoding/json"
	"github.com/opinionated/scraper-core/net"
	"github.com/opinionated/scraper-core/scraper"
	"github.com/opinionated/utils/log"
	"io/ioutil"
	"net/http"
	"time"
)

// TODO: handle connection errors
// TODO: scraping logic

type Client struct {
	IP        string
	ignoreBad map[string]struct{}
}

// Run runs the client in an infinite loop
func (c *Client) Run() {
	ticker := time.NewTicker(time.Duration(10) * time.Second)
	c.ignoreBad = make(map[string]struct{})
	defer c.printBad()

	for {

		// wait for the next itr
		<-ticker.C

		// go get the article
		req, err := Get(c.IP)
		if err != nil {
			log.Error(err)
			continue // go around to the next iteration
		}

		// don't reply to empty requests
		if netScraper.IsEmptyRequest(req) {
			log.Info("got empty request")
			continue
		}
		if _, contains := c.ignoreBad[req.URL]; contains {
			log.Warn("got bad article:", req.URL)
			continue
		}

		log.Info("got article", req.URL)

		// for now only use the NYT
		article := scraper.NYTArticle{}
		article.Link = req.URL

		err = scraper.ScrapeArticle(&article)
		if err != nil {
			c.ignoreBad[req.URL] = struct{}{} // add to ignore
			log.Error("could not scrape article", req.URL, ":", err)
			continue
		}

		if len(article.GetData()) == 0 {
			c.ignoreBad[req.URL] = struct{}{} // add to ignore
			log.Error("bad article body for url:", req.URL)
			continue
		}

		// send article back up
		result := netScraper.Response{URL: req.URL, Data: article.Data, Error: netScraper.ResponseOk}
		err = Post(c.IP, result)
		if err != nil {
			time.Sleep(time.Duration(10) * time.Second) // sleep incase the WIFI is just down

			err = Post(c.IP, result)
			if err != nil {
				log.Error(err)
			}
		}
	}
}

func (c Client) printBad() {
	for key, _ := range c.ignoreBad {
		log.Warn("ignored article:", key)
	}
}

// Get the server for an article to scrape.
func Get(target string) (netScraper.Request, error) {
	c := &http.Client{}

	// get next work unit
	resp, err := c.Get(target)

	if err != nil {
		return netScraper.EmptyRequest(), err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Error("oh nose, did not get OK status:", resp.StatusCode)
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
func Post(target string, done netScraper.Response) error {

	c := &http.Client{}

	fetchedJSON, err := json.Marshal(done)

	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST",
		target,
		bytes.NewReader(fetchedJSON))

	if err != nil {
		return err
	}

	req.Header.Set("content-type", "application/json")

	postResp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer postResp.Body.Close()

	return nil
}
