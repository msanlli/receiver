package pkg

import (
	"sync"
)

func Main() {
	// A wait group is implemented in order to wait for all goroutines to finish
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		startTCP()
		wg.Done()
	}()

	go func() {
		startUDP()
		wg.Done()
	}()

	wg.Wait()
}
