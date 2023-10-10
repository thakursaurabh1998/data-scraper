package main

import (
	"sync"

	_ "github.com/lib/pq"
	"github.com/thakursaurabh1998/scraper/scraper"
)

const MAX_CONC = 100

func main() {
	dbConn := scraper.GetDBConnection()
	defer dbConn.Close()

	var wg sync.WaitGroup

	// SAMPLE AWB PRD000085800
	// start := 85635
	start := 157462
	end := 157463

	sem := make(chan int, MAX_CONC)

	for generatedAWBNumber := range scraper.GenerateAWB(start, end) {
		sem <- 1
		wg.Add(1)
		go func(awbNumber string) {
			defer wg.Done()
			scraper.FetchOrderForAWB(awbNumber)
			<-sem
		}(generatedAWBNumber)
	}

	wg.Wait()
}
