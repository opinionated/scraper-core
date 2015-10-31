package scraper

import (
	"encoding/xml"
	"fmt"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
)

/**
 * Types for scraping articles from a web source
 * RSS feeds contain channels which contain articles
 *
 * Also provides UpdateRSS(RSS) and ScrapeArticle(Article)
 */

// Article holds scaping data(link) and parsed data (body, title, description)
// Gets body from an article link.
type Article interface {
	DoParse(*html.Tokenizer) error

	SetData(string)
	GetData() string

	GetLink() string
	GetDescription() string
	GetTitle() string
}

// RSS contains a link and a channel.
// Used when unmarshalling rss feeds.
type RSS interface {
	// TODO: add ptr and non-ptr access to these guys
	GetLink() string
	GetChannel() RSSChannel
}

// RSSChannel is basically an array of Articles.
// Can't return an array here because of how interfaces are set up.
type RSSChannel interface {
	GetArticle(int) Article
	GetNumArticles() int
	ClearArticles() bool
}

// UpdateRSS finds articles currently in the RSS feed.
// Clears old articles out of RSS feed before getting new ones.
// rss should be passed as an *RSS
func UpdateRSS(rss RSS) error {
	// clear out old articles so we don't double add
	rss.GetChannel().ClearArticles()

	// send request
	resp, err := http.Get(rss.GetLink())
	if err != nil {
		// TODO: error checking here
		fmt.Println("error getting RSS:", err)
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error reading body")
		// TODO: error handling
		return err
	}

	err = xml.Unmarshal(body, rss)
	if err != nil {
		fmt.Println("could not build RSS obj from rss request")
		return err
	}

	return nil
}

// ScrapeArticle fetches and parses the article.
// article should be provided as a *Article.
func ScrapeArticle(article Article) error {
	client := &http.Client{}

	// build request
	req, err := http.NewRequest("GET", article.GetLink(), nil) //create http request
	err = buildArticleHeader(req)
	if err != nil {
		fmt.Println("could not build article request")
		return err
	}

	//send http request
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println("error sending article request")
		return err
	}

	// TODO: check resp.Header to see if X-Article-Template is [full]

	// parse request
	parser := html.NewTokenizer(resp.Body)
	err = article.DoParse(parser) //parse the html body
	if err != nil {
		fmt.Println("error building article request")
		return err
	}
	return nil
}

// build headers for article request
func buildArticleHeader(req *http.Request) error {
	req.Header.Add("Referer", "https://www.google.com") //required to get past paywall

	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")                   //extra?
	req.Header.Add("Accept-Language", "en-US,en;q=0.5")                                                           //extra?
	req.Header.Add("Host", "www.wsj.com")                                                                         //extra?
	req.Header.Add("Cookie", "DJSESSION=country%3Dus%7C%7Ccontinent%3Dna%7C%7Cregion%3Dny%7C%7Ccity%3Dpoundtown") //messin with cookies
	return nil
}
