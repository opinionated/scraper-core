package server

import (
	"fmt"
	"github.com/opinionated/scraper-core/scraper"
	"net/http"
)

// TODO: make this an example?
func main() {

	b := func(name string) scraper.Article {
		return &scraper.WSJArticle{Link: name}
	}

	fmt.Println("startingServer")
	j := NewJefe()

	j.Add(b("one"))
	j.Add(b("two"))
	j.Add(b("three"))
	http.HandleFunc("/", j.Handle())
	http.ListenAndServe(":8080", nil)
}
