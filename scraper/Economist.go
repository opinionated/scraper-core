package scraper

import (
	"encoding/xml"
	"fmt"
	"golang.org/x/net/html"
	"strings"
)

// ECON new source types

type ECONArticle struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`

	// article body
	Data string
}

// TODO: add errors
func (article ECONArticle) GetLink() string        { return article.Link }
func (article ECONArticle) GetDescription() string { return article.Description }
func (article ECONArticle) GetTitle() string       { return article.Title }
func (article ECONArticle) GetData() string        { return article.Data }

// use ptrs for the next two because we want the article changed
func (article *ECONArticle) SetData(data string) { article.Data = data }

func (article *ECONArticle) DoParse(parser *html.Tokenizer) error {


// ENDS WITH div class content clearfix everywhere
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
					if attr.Key == "class" && attr.Val == "main-content" {
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

type ECONRSSChannel struct {
	XMLName  xml.Name     `xml:"channel"`
	Articles []ECONArticle `xml:"item"`
}

func (channel *ECONRSSChannel) GetArticle(slot int) Article {
	if slot >= channel.GetNumArticles() {
		// Check that the request doesn't go out of bounds
		// TODO: errors
		return nil
	}
	return &channel.Articles[slot]
}

func (channel *ECONRSSChannel) GetNumArticles() int {
	return len(channel.Articles)
}

type ECONRSS struct {
	XMLName xml.Name      `xml:"rss"`
	Channel ECONRSSChannel `xml:"channel"`
	RSSLink string
	// TODO: actually set string to the value of the link
}

func (rss *ECONRSS) GetLink() string { return "http://rss.ECONimes.com/services/xml/rss/ECON/HomePage.xml" }

func (rss *ECONRSS) GetChannel() RSSChannel {
	// return a pointer to the channel, interfaces implicitly have ptrs if they are there
	tmp := &rss.Channel
	return tmp
}

func (channel *ECONRSSChannel) ClearArticles() bool {
	channel.Articles = nil
	return true
}

// make sure all the structs implement the interfaces
var _ RSS = (*ECONRSS)(nil)
var _ RSSChannel = (*ECONRSSChannel)(nil)
var _ Article = (*ECONArticle)(nil)