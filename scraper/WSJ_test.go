package scraper_test

import (
	"bufio"
	"github.com/opinionated/scraper-core/fetcher"
	"github.com/opinionated/scraper-core/scraper"
	"golang.org/x/net/html"
	"os"
	"testing"
)

// TODO: make a comapare function
func TestTricky(t *testing.T) {
	file, err := os.Open("testData/WSJCarsonHtml.txt")
	defer file.Close()
	if err != nil {
		t.Errorf("error opening file:", err)
		return
	}

	fileScanner := bufio.NewReader(file)
	parser := html.NewTokenizer(fileScanner)

	article := &scraper.WSJArticle{}
	err = article.DoParse(parser)
	if err != nil {
		t.Errorf("error parsing:", err)
	}
}
