package main

import (
	"fmt"
	"github.com/lokken/family-football/gobbler"
	"log"
	"net/http"
	"sync"
)

func main() {
	downloadSchedules()
}

func downloadSchedules() {
	var wg sync.WaitGroup

	for i := 1; i <= 14; i++ {
		wg.Add(1)

		go func(week int) {
			defer wg.Done()

			url := fmt.Sprintf("https://www.espn.com/college-football/schedule/_/week/%d", week)

			resp, err := http.Get(url)
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()

			fmt.Printf("Downloading week %d ...\n", week)
			gobbler.SaveSchedule(week, resp.Body)
		}(i)
	}

	wg.Wait()
}
