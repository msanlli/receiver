//receiver v2.2.0

package main

import (
	pkg "receiver.com/m/pkg"
)

func main() {
	dir := []string{
		"./alert/alert.json",
		"./alert/alert.yaml",
		"./alert/alert.toml",
		"./data/data.json",
		"./data/data.yaml",
		"./data/data.toml",
	}

	pkg.Main(dir)
}
