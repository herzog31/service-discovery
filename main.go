package main

import (
	"flag"
	"log"
)

func main() {

	var dockerAPI string
	var webPort int64

	flag.StringVar(&dockerAPI, "api", "unix:///var/run/docker.sock", "Address of Docker API. Defaults to unix:///var/run/docker.sock")
	flag.Int64Var(&webPort, "port", 8080, "Port for the service discovery's API and web interface.")
	flag.Parse()

	d, err := NewDiscovery(dockerAPI, webPort)
	if err != nil {
		log.Fatal(err)
	}

	// Initial refresh
	d.refreshList()

	// Listen to Docker events
	go d.listen()

	// Serve API and web interface
	d.serveWeb()

	// TODO(mjb): JSON override omits
	// TODO(mjb): Notifications on Errors
	// TODO(mjb): Save logs
	// TODO(mjb): Persistence
	// TODO(mjb): Authentication
	// TODO(mjb): Docker container: Linux x64
	// TODO(mjb): Docker container: Linux ARM
}
