//subscrape.go
package scraper

import (
	"io"
	"fmt"
	// "net/url"
	// "net/http"		//http.Response
	"golang.org/x/net/html"
)

func GetSubResources(cj *CookieJar, body io.ReadCloser) {
	cj.Yo()
	// fmt.Println(body)
	// u, _ := url.Parse("https://google.com")
	// fmt.Println(cj.GetCookies(u)[0])
	parser := html.NewTokenizer(body)
	parseForResources(parser)
}

func parseForResources(parser *html.Tokenizer) {
	fmt.Println("starting parse")
	for {
		tt := parser.Next()
		switch tt {
			case html.ErrorToken:
				return
			case html.StartTagToken:
				tmp := parser.Token()
				if tmp.Data == "script" || tmp.Data == "img" {
					for _, a := range tmp.Attr {
					    if a.Key == "src" {
					        fmt.Println("Found", tmp.Data, "\t", a.Val)
					        break
					    }
					}
				}
		}	
	}
}