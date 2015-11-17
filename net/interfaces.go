// Package netScraper provides structures common to the
// client and the server.
package netScraper

const (
	// ResponseOk means file successfully parsed.
	ResponseOk = "OK"
	// ResponseBad when error parsing.
	ResponseBad = "BAD"
)

// Request is a URL to scrape.
type Request struct {
	URL string
}

// Response is a scraped URL response.
type Response struct {
	URL   string
	Data  string
	Error string
}

// EmptyRequest constructs and returns an empty request. Here
// so we can change the definition of an empty request later.
func EmptyRequest() Request {
	return Request{""}
}

// IsEmptyRequest checks if a request is empty. Here so we can
// change the definition of an empty request later.
func IsEmptyRequest(r Request) bool {
	return r.URL == ""
}
