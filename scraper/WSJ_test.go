package scraper_test

import (
	"bufio"
	"github.com/opinionated/scraper-core/scraper"
	"golang.org/x/net/html"
	"os"
	"testing"
)

// TODO: make a compare function
func TestWSJ1(t *testing.T) {
	file, err := os.Open("testData/WSJCarsonHtml.txt") // TODO: Add test file
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
