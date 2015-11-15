package main

import (
	"fmt"
	"github.com/opinionated/scraper-core/scraper"
	"net/http"
)

// need to make netRSS task to add scrape targets to a queue.
// the server pops the next thing off the stack when a request
// comes through.

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
