package main

import (
	"log"
)

func main() {

	d, err := NewDiscovery("tcp://192.168.178.27:4243", 8080)
	if err != nil {
		log.Fatal(err)
	}

	// Initial refresh
	d.refreshList()

	// Listen to Docker events
	go d.listen()

	// Serve API and web interface
	d.serveWeb()

	// TODO(mjb): Notifications on Errors
	// TODO(mjb): Save logs
	// TODO(mjb): Persistence
	// TODO(mjb): API Request to check if every container in environment is running or exited gracefully
	// TODO(mjb): Docker container: Linux x64
	// TODO(mjb): Docker container: Linux ARM
}
