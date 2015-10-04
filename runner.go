package main

import (
	"fmt"
	"github.com/opinionated/scraper-core/fetcher"
	"github.com/opinionated/scraper-core/scraper"
	"io/ioutil"
	"net/http"
)

func main() {
	// do a simple http fetch:
	resp, err := http.Get("http://www.wsj.com/xml/rss/3_7041.xml")
	if err != nil {
		fmt.Println("OH NOSE: got an error when trying to fetch the datz:", err)
		return
	}

	// make sure the body gets closed laster
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Oh nose: error reading body:", err)
		return
	}
	// note need to use ptr here because things will be changing
	rss := &scraper.WSJRSS{}
	err = fetcher.GetStories(rss, body)
	if err != nil {
		fmt.Println("oh nose, error working with body")
		return
	}
	err = fetcher.DoGetArticle(rss.GetChannel().GetArticle(0))
	fmt.Println("article body is:", rss.GetChannel().GetArticle(0).GetData())
}
