package generator

import (
	"fmt"
	"sync"
	"time"
)

// Simulator function to generate and send API requests
func Simulator(apiCount int, interval time.Duration, collectorURL string) {
	var wg sync.WaitGroup
	startTime := time.Now()

	for i := 0; i < apiCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			request := GetRandomRequest() 

			// Send request to collector
			err := request.SendRequest(collectorURL)
			if err != nil {
				fmt.Println("Request error:", err)
			}
		}()
		time.Sleep(interval) 
	}

	wg.Wait() // Wait for all goroutines to finish
	// fmt.Printf("Total time taken: %v\n", time.Since(startTime))
	fmt.Printf("Total time taken: %.2f seconds\n", time.Since(startTime).Seconds())

}
