package fetcher

import (
	"encoding/xml"
	"fmt"
	"golang.org/x/net/html"
	"strings"
)

// WSJ new source types

type WSJArticle struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`

	// article body
	Data string
}

// TODO: add errors
func (article WSJArticle) GetLink() string        { return article.Link }
func (article WSJArticle) GetDescription() string { return article.Description }
func (article WSJArticle) GetTitle() string       { return article.Title }
func (article WSJArticle) GetData() string        { return article.Data }

// use ptrs for the next two because we want the article changed
func (article *WSJArticle) SetData(data string) { article.Data = data }

func (article *WSJArticle) DoParse(parser *html.Tokenizer) error {

	// find the start of the article
	// starts at the top of the html body, ends at the article tag
articleTagLoop:
	for {
		token := parser.Next()

		switch {
		case token == html.ErrorToken:
			fmt.Println("OH NOSE!!!! ERROR before we hit the end")
			return nil
		case token == html.StartTagToken:
			tmp := parser.Token()

			isStartArticle := tmp.Data == "article"
			if isStartArticle {
				break articleTagLoop
			}
		}
	}

	// find the article header, which has author, time etc
	// starts at the article tag, ends at the article header
	// TODO: get author info and such here
articleStartLoop:
	for {
		token := parser.Next()

		switch {
		case token == html.ErrorToken:
			return nil
		case token == html.StartTagToken:
			tmp := parser.Token()

			isStartArticleBody := tmp.Data == "div"
			// loop until we are at the first paragraph of the article body
			if isStartArticleBody {
				isStartArticleBody = false
				for _, attr := range tmp.Attr {
					if attr.Key == "class" && attr.Val == "clearfix byline-wrap" {
						isStartArticleBody = true
						break
					}
				}
				if isStartArticleBody {
					break articleStartLoop
				}
			}
		}
	}

	// find the start of the article
	// starts at the end of the article header, ends at the first article paragraph
articleBodyStartLoop:
	for {
		token := parser.Next()
		switch {
		case token == html.ErrorToken:
			return nil
		case token == html.StartTagToken:
			tmp := parser.Token()
			isStartArticleBody := tmp.Data == "p"
			if isStartArticleBody {
				break articleBodyStartLoop
			}
		}
	}

	// pull the article out of the html
	// starts at first paragraph, returns at the end of the article
	isInParagraph := true // true because we start inside the first paragraph
	depth := 1            // one because this loop starts at first paragraph
	for {
		token := parser.Next()
		switch {
		case token == html.ErrorToken:
			fmt.Println("hit err, depth is:", depth)
			return nil
		case token == html.StartTagToken:
			depth++
			tmp := parser.Token()

			isParagraph := tmp.Data == "p"
			if isParagraph {
				// start of a new paragraph
				if depth != 1 {
					fmt.Println("ERROR: hit new paragraph while depth != 0")
				}
				if isInParagraph {
					fmt.Println("ERROR: hit unexpected new paragraph tag while in paragraph")
				}
				isInParagraph = true
			}

			// text can have embeded links
			isLink := tmp.Data == "a"
			if isLink {
				if !isInParagraph {
					fmt.Println("ERROR: hit unexpected link outside of a paragraph")
					continue
				}

				// if we are in a paragraph, append the link name
				parser.Next()
				tmp = parser.Token()
				newBody := article.GetData() + tmp.Data
				article.SetData(newBody)
			}
		case token == html.EndTagToken:
			depth--
			tmp := parser.Token().Data
			if depth == -1 {
				// done with article when we are at a higher level than it
				return nil
			}

			if tmp == "p" {
				// add a paragraph and trim the space
				article.SetData(strings.TrimSpace(article.GetData() + "\n"))
				isInParagraph = false
			}

		default:
			if !isInParagraph {
				// if not inside a text paragraph, continue on
				continue
			}

			// get the paragraph text and append it to the article body
			// TODO: look into using a string builder instead of adding things on
			tmp := parser.Token()
			newBody := article.GetData()
			// add a space on the left just in case there is a comment or something
			newBody = newBody + strings.TrimSpace(tmp.Data) + " "

			article.SetData(newBody)
		}

	}
	return nil
}

type WSJRSSChannel struct {
	XMLName  xml.Name     `xml:"channel"`
	Articles []WSJArticle `xml:"item"`
}

func (channel *WSJRSSChannel) GetArticle(slot int) Article {
	if slot >= channel.GetNumArticles() {
		// Check that the request doesn't go out of bounds
		// TODO: errors
		return nil
	}
	return &channel.Articles[slot]
}

func (channel *WSJRSSChannel) GetNumArticles() int {
	return len(channel.Articles)
}

type WSJRSS struct {
	XMLName xml.Name      `xml:"rss"`
	Channel WSJRSSChannel `xml:"channel"`
	RSSLink string
	// TODO: actually set string to the value of the link
}

func (rss *WSJRSS) GetLink() string { return "http://www.wsj.com/xml/rss/3_7041.xml" }

func (rss *WSJRSS) GetChannel() RSSChannel {
	// return a pointer to the channel, interfaces implicitly have ptrs if they are there
	tmp := &rss.Channel
	return tmp
}

// make sure all the structs implement the interfaces
var _ RSS = (*WSJRSS)(nil)
var _ RSSChannel = (*WSJRSSChannel)(nil)
var _ Article = (*WSJArticle)(nil)
