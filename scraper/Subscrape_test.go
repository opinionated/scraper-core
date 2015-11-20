// Subscrape_test.go
package scraper_test

import (
	"fmt"
	"github.com/opinionated/scraper-core/scraper"
	"net/http" //http.Get
	"net/url"  //url.Parse
	"testing"

)

func TestSub1(t *testing.T) {

	client := &http.Client{}
	gurl := "https://reddit.com"
	req, err := http.NewRequest("GET", gurl, nil) //create http request
	if err != nil {
		fmt.Println("could not build article request")
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("could make net request")
		fmt.Println(err)
		return
	}

	cj := scraper.NewCookieJar()
	u, _ := url.Parse(gurl)
	cj.SetCookiesFromHeader(u, resp.Header)

	scraper.GetSubResources(cj, resp.Body)
}