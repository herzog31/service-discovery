package main

import (
	"errors"
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"strconv"
	"strings"
)

type Discovery struct {
	dockerAPI      string
	listener       chan *docker.APIEvents
	containers     []docker.APIContainers
	containersFull map[string]*docker.Container
	client         *docker.Client
}

func NewDiscovery(dockerAPI string) (*Discovery, error) {
	d := new(Discovery)
	d.dockerAPI = dockerAPI
	d.listener = make(chan *docker.APIEvents)
	d.containersFull = make(map[string]*docker.Container)
	client, err := docker.NewClient(d.dockerAPI)
	if err != nil {
		return nil, err
	}
	d.client = client
	return d, nil
}

func (d *Discovery) listen() {
	d.client.AddEventListener(d.listener)
	defer d.client.RemoveEventListener(d.listener)
	for {
		event := <-d.listener
		go d.handleEvent(event)
	}
}

func (d *Discovery) refreshList() error {
	options := docker.ListContainersOptions{
		All: true,
	}
	containers, err := d.client.ListContainers(options)
	if err != nil {
		return err
	}
	d.containers = containers
	for _, container := range containers {
		id := container.ID
		full, err := d.client.InspectContainer(id)
		if err != nil {
			return err
		}
		d.containersFull[strings.TrimPrefix(full.Name, "/")] = full
	}
	return nil
}

func (d *Discovery) handleEvent(event *docker.APIEvents) error {
	fmt.Printf("Incoming Event: %v (%+v)\n", event.Status, event)
	err := d.refreshList()
	if err != nil {
		return err
		fmt.Printf("handleEvent Error: %+v", err)
	}
	return nil
}

func (d *Discovery) GetPortMappings(name string) (map[docker.Port][]docker.PortBinding, error) {
	container, ok := d.containersFull[name]
	if !ok {
		return nil, errors.New("Container not found!")
	}
	return container.NetworkSettings.Ports, nil
}

func (d *Discovery) GetPortMapping(name string, port Port) (Port, error) {
	mappings, err := d.GetPortMappings(name)
	if err != nil {
		return Port{}, err
	}

	internal := port.ToDockerPort()
	mapping, ok := mappings[internal]
	if !ok {
		return Port{}, errors.New(fmt.Sprintf("No port mapping available for internal port %s", port))
	}

	firstMapping := mapping[0]
	parsed, _ := strconv.ParseInt(firstMapping.HostPort, 10, 64)
	return Port{
		port:     parsed,
		protocol: port.protocol,
	}, nil
}
