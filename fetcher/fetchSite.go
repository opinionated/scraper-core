package fetcher

import (
	"encoding/xml"
	"fmt"
	"golang.org/x/net/html"
	"net/http"
)

// Parse the wsj opinion rss feed into useable links.
// Returns an array of articles and an error
func GetStories(rss RSS, body []byte) error {
	err := xml.Unmarshal(body, rss)
	if err != nil {
		fmt.Printf("err:", err)
		return err
	}

	for i := 0; i < rss.GetChannel().GetNumArticles(); i++ {
		article := rss.GetChannel().GetArticle(i)
		fmt.Println("title:", article.GetTitle(), "\tdescr:", article.GetDescription())
	}

	return nil
}

// Request a page containing the article linked to
func DoGetArticle(article Article) error {
	client := &http.Client{}

	req, err := http.NewRequest("GET", article.GetLink(), nil)
	req.Header.Add("Referer", "https://www.google.com")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("oh nose, err with get article http request:", err)
		return err
	}

	defer resp.Body.Close()
	parser := html.NewTokenizer(resp.Body)
	tmp := article.(*WSJArticle)
	tmp.DoParse(parser)
	return err
}
