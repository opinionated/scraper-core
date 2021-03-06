package scraper

import (
	"encoding/xml"
	"github.com/opinionated/utils/log"
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
		log.Error("error getting RSS:", err)
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("error reading RSS body:", err)
		// TODO: error handling
		return err
	}

	err = xml.Unmarshal(body, rss)
	if err != nil {
		log.Error("could not build RSS obj from rss request:", err)
		return err
	}

	return nil
}

// ScrapeArticle fetches and parses the article.
// article should be provided as a *Article.
func ScrapeArticle(article Article) error {
	cookies := NewCookieJar()
	client := &http.Client{Jar: cookies}

	// build request
	req, err := http.NewRequest("GET", article.GetLink(), nil) //create http request
	err = buildArticleHeader(req)
	if err != nil {
		log.Error("could not build article request:", err)
		return err
	}

	//send http request
	resp, err := client.Do(req)
	if err != nil {
		log.Error("error sending article request:", err)
		return err
	}
	defer resp.Body.Close()

	// TODO: check resp.Header to see if X-Article-Template is [full]

	// parse request
	parser := html.NewTokenizer(resp.Body)
	err = article.DoParse(parser) //parse the html body
	if err != nil {
		log.Error("error building article request:", err)
		return err
	}
	return nil
}

// build headers for article request
func buildArticleHeader(req *http.Request) error {
	req.Header.Add("Referer", "https://www.google.com") //required to get past paywall

	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8") //extra?
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/41.0.2228.0 Safari/537.36")
	req.Header.Add("Accept-Language", "en-US,en;q=0.5")                                                           //extra?
	req.Header.Add("Host", "www.nytimes.com")                                                                     //extra?
	req.Header.Add("Cookie", "DJSESSION=country%3Dus%7C%7Ccontinent%3Dna%7C%7Cregion%3Dny%7C%7Ccity%3Dpoundtown") //messin with cookies
	return nil
}
