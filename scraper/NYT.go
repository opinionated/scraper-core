package scraper

import (
	"encoding/xml"
	"fmt"
	"golang.org/x/net/html"
	"strings"
	"unicode"
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
			return fmt.Errorf("problem moving article %s to open tag", article.GetTitle())
		case token == html.StartTagToken:
			tmp := parser.Token()
			isStartArticle := tmp.Data == "p"
			if isStartArticle {
				for _, attr := range tmp.Attr {
					if attr.Key == "class" && attr.Val == "story-body-text story-content" {
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
			return fmt.Errorf("problem scraping article %s", article.GetTitle())
		case token == html.StartTagToken:
			tmp := parser.Token()
			isEndArticle := tmp.Data == "footer"
			if isEndArticle {
				for _, attr := range tmp.Attr {
					if attr.Key == "class" && attr.Val == "story-footer story-content" {
						break articleClosingTagLoop
					}
				}
			}

			if tmp.Data == "p" {
				for _, attr := range tmp.Attr {
					if attr.Key == "class" && strings.Contains(attr.Val, "story-body-text") {
						isInParagraph = true
					}
				}
				if isInParagraph {
					continue
				}
			}

			// is a link
			if tmp.Data == "a" {
				shouldSkip := false
				for _, attr := range tmp.Attr {
					if attr.Key == "class" && strings.Contains(attr.Val, "visually-hidden") {
						shouldSkip = true
					}
				}

				if shouldSkip {
					continue
				}

				parser.Next()
				tmp = parser.Token()
				newBody := strings.TrimSpace(article.GetData()) + " " + strings.TrimSpace(tmp.Data) + " "
				article.SetData(newBody)
				isInParagraph = true
			}

		case token == html.EndTagToken:
			tmp := parser.Token()
			if tmp.Data == "p" {
				isInParagraph = false
			}

		default:
			if !isInParagraph {
				continue
			}
			tmp := parser.Token()

			newBody := article.GetData()
			// add a space on the left just in case there is a comment or something
			if unicode.IsPunct(rune(tmp.Data[0])) {
				newBody = strings.TrimSpace(newBody)
			}
			newBody = newBody + strings.TrimSpace(tmp.Data)
			article.SetData(newBody)
			isInParagraph = false
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

func (rss *NYTRSS) GetLink() string {
	return "http://topics.nytimes.com/top/opinion/editorialsandoped/editorials/index.html?rss=1"
}

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
var _ Article = (*NYTArticle)(nil)
