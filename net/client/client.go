package main

import (
	"bytes"
	"fmt"
	"github.com/opinionated/scraper-core/net"
	"io/ioutil"
	"net/http"
	//"strings"
	"encoding/json"
)

func main() {

	c := &http.Client{}

	// get next work unit
	resp, err := c.Get("http://localhost:8080/")
	defer resp.Body.Close()
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 200 {
		fmt.Println("oh nose, did not get OK status:", resp.StatusCode)
		panic("bad status code")
	}

	js, err := ioutil.ReadAll(resp.Body)
	fmt.Println("response body is:", string(js))
	toDo := netScraper.Request{}

	err = json.Unmarshal(js, &toDo)
	if err != nil {
		panic(err)
	}

	fetched := netScraper.Response{toDo.Url, "data"}
	fetchedJson, err := json.Marshal(fetched)
	if err != nil {
		panic(err)
	}
	req, err := http.NewRequest("POST",
		"http://localhost:8080/",
		bytes.NewReader(fetchedJson))

	if err != nil {
		panic(err)
	}
	req.Header.Set("content-type", "application/json")

	postResp, err := c.Do(req)
	defer postResp.Body.Close()
	if err != nil {
		panic(err)
	}

}
