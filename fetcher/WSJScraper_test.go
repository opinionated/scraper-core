package fetcher_test

import (
	"bufio"
	"golang.org/x/net/html"
	"os"
	"scraper/fetcher"
	"testing"
)

// TODO: make a comapre function
func TestTricky(t *testing.T) {
	file, err := os.Open("testData/WSJCarsonHtml.txt")
	defer file.Close()
	if err != nil {
		t.Errorf("error opening file:", err)
		return
	}

	fileScanner := bufio.NewReader(file)
	parser := html.NewTokenizer(fileScanner)

	article := &fetcher.WSJArticle{}
	err = article.DoParse(parser)
	if err != nil {
		t.Errorf("error parsing:", err)
	}
}
