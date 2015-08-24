package main

import (
	"os"
)

func main() {

	/* d := &Discovery{
		endpoint: "tcp://192.168.178.27:4243",
		listener: make(chan *docker.APIEvents),
	} */

	d, err := NewDiscovery("tcp://192.168.178.27:4243")
	if err != nil {
		os.Exit(1)
	}
	d.listen()

	// TODO(mjb): Get Port Mapping
	// TODO(mjb): Refresh Data on Event
	// TODO(mjb): Settings
	// TODO(mjb): Webserver + API
	// TODO(mjb): Save logs
	// TODO(mjb): Persistence

	//eventHandler(events)

	/* imgs, _ := client.ListImages(docker.ListImagesOptions{All: false})
	for _, img := range imgs {
		fmt.Println("ID: ", img.ID)
		fmt.Println("RepoTags: ", img.RepoTags)
		fmt.Println("Created: ", img.Created)
		fmt.Println("Size: ", img.Size)
		fmt.Println("VirtualSize: ", img.VirtualSize)
		fmt.Println("ParentId: ", img.ParentID)
	} */
}
