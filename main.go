package main

import (
	"os"
)

func main() {

	d, err := NewDiscovery("tcp://192.168.178.27:4243")
	if err != nil {
		os.Exit(1)
	}
	d.listen()

	// TODO(mjb): Settings
	// TODO(mjb): Webserver + API
	// TODO(mjb): Save logs
	// TODO(mjb): Persistence
}
