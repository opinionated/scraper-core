package netScraper_test

import (
	"fmt"
	//	"github.com/opinionated/scraper-core/net"
	"github.com/opinionated/scraper-core/net/client"
	"github.com/opinionated/scraper-core/net/server"
	"github.com/opinionated/scraper-core/scraper"
	//	"net/http"
	"testing"
	"time"
)

func doFullServer() *server.Jefe {
	j := server.NewJefe()

	// get the jefe going
	go StartServer(&j)

	// make the scheduler loop quick
	j.SetCycleTime(1)
	j.Start()

	// build RSS and add it
	rss := server.CreateSchedulableRSS(&scraper.WSJRSS{}, 10, &j)
	j.AddSchedulable(rss)

	return &j
}

func TestIntegrationA(t *testing.T) {
	j := doFullServer()

	c := client.Client{}
	go c.Run()
	time.Sleep(time.Duration(100) * time.Second)

	fmt.Println("going to stop j")
	j.Stop()
}
