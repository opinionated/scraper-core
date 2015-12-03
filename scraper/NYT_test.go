package scraper_test

import (
	"bufio"
	"fmt"
	"github.com/opinionated/scraper-core/scraper"
	"golang.org/x/net/html"
	"os"
	"testing"
)

// TODO: make a compare function
func TestNYT1(t *testing.T) {
	err := CompareBodies("testData/NYTPutinHTML.txt", "testData/NYTPutinBody.txt", &scraper.NYTArticle{})
	if err != nil {
		t.Errorf("%s", err)
	}
}

func TestCompareBodies(t *testing.T) {
	err := CompareBodies("testData/NYTPutinHTML.txt", "testData/NYTPutinHTML.txt", &scraper.NYTArticle{})
	if err != nil {
		t.Errorf("%s", err)
	}
}

// CompareBodies will take two strings, one HTML path and one path to a
// pre-parsed file and compare the two using the parse method from article.
// If they are the same it returns nil, otherwise will return an error.
func CompareBodies(HTMLarticle string, ExpectedArticle string, article scraper.Article) error {
	file, err := os.Open(HTMLarticle) // TODO: Add test file
	defer file.Close()
	if err != nil {
		return fmt.Errorf("error opening HTML file %s", err)
	}

	fileScanner := bufio.NewReader(file)
	parser := html.NewTokenizer(fileScanner)

	err = article.DoParse(parser)
	if err != nil {
		return fmt.Errorf("error parsing: %s", err)
	}

	fileCompare, err := os.Open(ExpectedArticle)
	defer fileCompare.Close()
	if err != nil {
		return fmt.Errorf("error opening body file:%s", err)
	}

	CompareFile := bufio.NewScanner(fileCompare)
	fullText := ""
	for CompareFile.Scan() {
		fullText += CompareFile.Text()
	}
	if fullText != article.GetData() {
		return fmt.Errorf("Expected: \n%s\n Received: \n%s\n", fullText, article.GetData())
	}

	return nil
}
