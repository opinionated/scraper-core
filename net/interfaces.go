package netScraper

type Request struct {
	Url string
}

type Response struct {
	Url  string
	Data string
}

func EmptyRequest() Request {
	return Request{""}
}

func IsEmptyRequest(r Request) bool {
	return r.Url == ""
}
