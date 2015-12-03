package scraper

import (
	"encoding/xml"
	"fmt"
	"golang.org/x/net/html"
	"strings"
)

// NYT new source types

type NYTArticle struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`

	// article body
	Data string
}

// TODO: add errors
func (article NYTArticle) GetLink() string        { return article.Link }
func (article NYTArticle) GetDescription() string { return article.Description }
func (article NYTArticle) GetTitle() string       { return article.Title }
func (article NYTArticle) GetData() string        { return article.Data }

// use ptrs for the next two because we want the article changed
func (article *NYTArticle) SetData(data string) { article.Data = data }

func (article *NYTArticle) DoParse(parser *html.Tokenizer) error {

articleOpeningTagLoop:
	for {
		token := parser.Next()

		switch {
		case token == html.ErrorToken:
			fmt.Println("Prollem")
			return nil
		case token == html.StartTagToken:
			tmp := parser.Token()
			isStartArticle := tmp.Data == "p"
			if isStartArticle {
				for _, attr := range tmp.Attr {
					if attr.Key == "class" && attr.Val == "story-body-text story-content" {
						fmt.Println("Found it, bitch")
						break articleOpeningTagLoop
					}
				}
			}
		}
	}

	isInParagraph := true
articleClosingTagLoop:
	for {
		token := parser.Next()
		switch {
		case token == html.ErrorToken:
			fmt.Println("Prollem")
			return nil
		case token == html.StartTagToken:
			tmp := parser.Token()
			isEndArticle := tmp.Data == "footer"
			if isEndArticle {
				for _, attr := range tmp.Attr {
					if attr.Key == "class" && attr.Val == "story-footer story-content" {
						fmt.Println("Hit end")
						break articleClosingTagLoop
					}
				}
			}
			isInParagraph = true
		default:
			if !isInParagraph {
				continue
			}
			tmp := parser.Token()

			newBody := article.GetData()
			// add a space on the left just in case there is a comment or something
			newBody = newBody + strings.TrimSpace(tmp.Data)
			article.SetData(newBody)
			isInParagraph = false
			//fmt.Println("Next p", newBody)
		}
	}
	fmt.Println(article.GetData())
	return nil
}

type NYTRSSChannel struct {
	XMLName  xml.Name     `xml:"channel"`
	Articles []NYTArticle `xml:"item"`
}

func (channel *NYTRSSChannel) GetArticle(slot int) Article {
	if slot >= channel.GetNumArticles() {
		// Check that the request doesn't go out of bounds
		// TODO: errors
		return nil
	}
	return &channel.Articles[slot]
}

func (channel *NYTRSSChannel) GetNumArticles() int {
	return len(channel.Articles)
}

type NYTRSS struct {
	XMLName xml.Name      `xml:"rss"`
	Channel NYTRSSChannel `xml:"channel"`
	RSSLink string
	// TODO: actually set string to the value of the link
}

func (rss *NYTRSS) GetLink() string { return "http://topics.nytimes.com/top/opinion/editorialsandoped/editorials/index.html?rss=1" }

func (rss *NYTRSS) GetChannel() RSSChannel {
	// return a pointer to the channel, interfaces implicitly have ptrs if they are there
	tmp := &rss.Channel
	return tmp
}

func (channel *NYTRSSChannel) ClearArticles() bool {
	channel.Articles = nil
	return true
}

// make sure all the structs implement the interfaces
var _ RSS = (*NYTRSS)(nil)
var _ RSSChannel = (*NYTRSSChannel)(nil)
var _ Article = (*WSJArticle)(nil)
