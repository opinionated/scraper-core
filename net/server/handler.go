package server

import (
	"encoding/json"
	"github.com/opinionated/scraper-core/scraper"
	"github.com/opinionated/utils/log"
	"io/ioutil"
	"strings"
)

// handles scraped articles
// TODO: think about where this should be
func handleScrapedArticle(article scraper.Article) {
	if err := scraper.CheckFile(article.GetData()); err != nil {
		log.Warn("when checking article", article.GetTitle(), "got err:", err)
	}
	if err := storeArticle(article); err != nil {
		log.Error("failed to write article", article.GetTitle(), ":", err)
		return
	}
}

// write the provided article to the storage
func storeArticle(article scraper.Article) error {
	jsonStr, err := json.Marshal(article)
	if err != nil {
		return err
	}

	// take all spaces out of title
	// TODO: think about cleaning this up a little more
	fileName := strings.Replace(article.GetTitle(), " ", "", -1)
	path := "opinionatedData/" + fileName + ".json"

	err = ioutil.WriteFile(path, jsonStr, 0644)
	if err != nil {
		return err
	}

	log.Info("wrote article:", article.GetTitle(), "to location:", path)
	return nil
}
