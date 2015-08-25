package main

import (
	"os"
)

func main() {

	d, err := NewDiscovery("tcp://192.168.178.27:4243", 8080)
	if err != nil {
		os.Exit(1)
	}
	d.refreshList()
	go d.listen()
	d.serveWeb()

	// TODO(mjb): Settings
	// TODO(mjb): Save logs
	// TODO(mjb): Persistence
}
