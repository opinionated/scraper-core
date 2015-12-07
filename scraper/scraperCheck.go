package scraper

import (
	"fmt"
)

// CheckFile looks at a text string and recognizes illegal characters. Then it detemines
// likelihood of incorrect scraping.
func CheckFile(FileToCheck string) error {
	numPotError := 0
	para := 0

	for _, c := range FileToCheck {
		// Checks for multiple paragraphs in a row,
		// and backslashes in text body.
		if c == '\n' && para == 0 {
			para = 1
		} else if c == '\n' && para == 1 {
			numPotError += 1
		}
		if c == '\\' || c == '/' {
			numPotError += 5
		}

		charAsc := c
		// Looks at ASCII value of each character, if it is not English alphabet
		// it adds 5 to the error count.
		if charAsc > 122 || charAsc < 9 || charAsc == 11 || charAsc == 12 || (8 < charAsc && charAsc < 32) {
			numPotError += 5
		}
	}

	if numPotError == 0 {
		return nil
	} else if numPotError > 0 && numPotError < 11 {
		return fmt.Errorf("%d errors found, potentially safe article.", numPotError)
	} else {
		return fmt.Errorf("WARNING: MANY ERRORS FOUND, CHECK ARTICLE MANUALLY.")
	}
}
