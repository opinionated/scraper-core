package fetcher

import (
	"golang.org/x/net/html"
)

/**
 * Types for scraping articles from a web source
 * Most news sites seem to have rss feeds for articles
 * RSS feeds contain channels which contain articles
 * Use these interfaces to add a new news source
 */
// Generic article type
// Gets body from an article link
type Article interface {
	DoParse(*html.Tokenizer) error

	SetData(string)
	GetData() string

	GetLink() string
	GetDescription() string
	GetTitle() string
}

// Generic RSS feed
// RSS feeds have channels
type RSS interface {
	// TODO: add ptr and non-ptr access to these guys
	GetLink() string
	GetChannel() RSSChannel
}

// Provides access to an RSS's channel, which contains the articles
// Can't return an array here because of how interfaces are set up
type RSSChannel interface {
	GetArticle(int) Article
	GetNumArticles() int
}
