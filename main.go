//receiver v2.1.3

package main

import (
	"os"
	"os/signal"
	"syscall"

	pkg "receiver.com/m/pkg"
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

func main() {
	record, err := pkg.NewRecord([]string{
		"./data/data.json",
		"./data/data.yaml",
		"./data/data.toml",
		"./alert/alert.json",
		"./alert/alert.yaml",
		"./alert/alert.toml",
	})
	if err != nil {
		panic(err)
	}
	defer record.CloseAll()

	closeHandler(record.CloseAll)

	pkg.Main()
}
