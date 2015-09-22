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

	// Gather logs
	go d.gatherLogs()

	// Serve API and web interface
	d.serveWeb()

	// TODO(mjb): Add template folder to releases
	// TODO(mjb): Prevent error if no log exists, show that log is empty
	// TODO(mjb): JSON override omits
	// TODO(mjb): Persistence (redis, settings, logs)
	// TODO(mjb): Authentication (in settings)
	// TODO(mjb): API Documentation (separate tab, link to github)
	// TODO(mjb): CD Docker container: Linux x64
	// TODO(mjb): CD Docker container: Linux ARM
}
