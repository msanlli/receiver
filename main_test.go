package main_test

import (
	"testing"

	main "receiver.com/m"
)

func TestMain(t *testing.T) {
	go main.Main()

}
