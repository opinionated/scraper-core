package scraper_test

import (
	"bufio"
	"testing"
	"github.com/opinionated/scraper-core/scraper"
	"os"
)

func ExpectFail(input string, t *testing.T) {
	err := scraper.CheckFile(input)
	if err == nil {
		t.Errorf("%s Failed on input %s", err, input)
	}
}

func TestScrapeCheckFail(t *testing.T) {
	ExpectFail("\\\\\\", t)
	ExpectFail("\n\n\n\n\n\n\n\n\n\n\n\n\n\n\n", t)
	ExpectFail("\x005", t)
	ExpectFail("\x19", t)
	ExpectFail("\x0B", t)
	ExpectFail("\x0C", t)
	ExpectFail("\x7C", t)
	ExpectFail("\x1F", t)
}

func TestScrape(t *testing.T) {
	file, err := os.Open("testData/NYTPutinBody.txt")
	if err != nil {
		t.Errorf("error opening file %s", err)
	}
	defer file.Close()

	scanned := bufio.NewScanner(file)
	stringArticle := ""
	for scanned.Scan() {
		stringArticle += scanned.Text()
	}
	err = scraper.CheckFile(stringArticle)
	if err != nil {
		t.Errorf("%s",err)
	}
}