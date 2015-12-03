package netScraper_test

import (
	"github.com/opinionated/scraper-core/net/client"
	"github.com/opinionated/scraper-core/net/server"
	"github.com/opinionated/scraper-core/scraper"
	"github.com/opinionated/utils/log"
	"testing"
	"time"
)

func doFullServer() *server.Jefe {
	s := server.NewScrapeServer()
	j := s.GetJefe()

	// get the jefe going
	go StartServer(s)

	// make the scheduler loop quick
	j.SetCycleTime(1)
	j.Start()

	// build RSS and add it
	rss := server.CreateSchedulableRSS(&scraper.WSJRSS{}, 10, j)
	j.AddSchedulable(rss)

	return j
}

func TestIntegrationA(t *testing.T) {
	log.InitStd()

	j := doFullServer()
	c := client.Client{}
	go c.Run()
	time.Sleep(time.Duration(100) * time.Second)

	log.InitStd()
	log.Info("going to stop j")
	j.Stop()
}
