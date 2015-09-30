package fetcher

import (
	"encoding/xml"
	"fmt"
	"golang.org/x/net/html"
	"io/ioutil"
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
	rss := &WSJRSS{}
	err = GetStories(rss, body)
	if err != nil {
		fmt.Println("oh nose, error working with body")
		return
	}
	err = DoGetArticle(rss.GetChannel().GetArticle(0))
	fmt.Println("article body is:", rss.GetChannel().GetArticle(0).GetData())
}
