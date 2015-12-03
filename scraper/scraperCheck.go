package scraper

import (
	"fmt"
)

func CheckFile(FileToCheck string) error {
	numPotError := 0	
	para := 0
	for _,c := range FileToCheck {
		if c == '\n' && para == 0 {
			para = 1
		} else if c == '\n' && para == 1 {
			numPotError += 1
		}
		if c == '\\' || c == '/' {
			numPotError += 5
		}
		charAsc := c
		if charAsc > 122 || charAsc < 9 || charAsc == 11 || charAsc == 12 || (8 < charAsc && charAsc < 32)  {
			numPotError += 5
		}
	}
	if numPotError == 0 {
		fmt.Println("No potential errors.")
		return nil
	} else if numPotError > 0 && numPotError < 11 {
		return fmt.Errorf("%d errors found, potentially safe article.", numPotError)
	} else {
		return fmt.Errorf("WARNING: MANY ERRORS FOUND, CHECK ARTICLE MANUALLY.")
	}
}