// cookie_test.go
package scraper_test

import (
	"fmt"
	"net/http"				//http.Get
	"net/url"
	"github.com/opinionated/scraper-core/scraper"
	"testing"
)

func TestCookie1(t *testing.T) {
	t.Skip("skipping cookie test for now")

	cj := scraper.NewCookieJar()

	client := &http.Client{}
	urlstring := "https://google.com"
	u, _ := url.Parse(urlstring) 
	req, _ := http.NewRequest("GET", urlstring, nil) //create http request
	res, err := client.Do(req)
	defer res.Body.Close()
	if err != nil {
		fmt.Println("error sending article request")
	}

	cj.SetCookiesFromHeader(u, res.Header)
	u2, _ := url.Parse("https://docs.google.com/penis")
	cj.SetCookies(u, nil)
	fmt.Println(cj.GetCookies(u2))
}