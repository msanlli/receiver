package pkg

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func closeHandler(cleanupFunc func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanupFunc()
		os.Exit(0)
	}()
}

func Main(dir []string) {
	var alert, _ = NewAlertRecord(dir)
	var data, _ = NewDataRecord(dir)

	start, err := NewRecord(dir)
	if err != nil {
		println("Error opening files:", err)
		return
	}

	closeHandler(start.CloseAll)

	// A wait group is implemented in order to wait for all goroutines to finish
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		startTCP(alert, data)
		wg.Done()
	}()

	go func() {
		startUDP(alert, data)
		wg.Done()
	}()

	wg.Wait()
}
