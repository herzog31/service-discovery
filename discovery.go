package main

import (
	"errors"
	"fmt"
	"github.com/fsouza/go-dockerclient"
)

type Discovery struct {
	dockerAPI  string
	listener   chan *docker.APIEvents
	containers []docker.APIContainers
	client     *docker.Client
}

func NewDiscovery(dockerAPI string) (*Discovery, error) {
	d := new(Discovery)
	d.dockerAPI = dockerAPI
	d.listener = make(chan *docker.APIEvents)
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
	return nil
}

func (d *Discovery) handleEvent(event *docker.APIEvents) error {
	fmt.Printf("%+v\n", event)
	err := d.refreshList()
	if err != nil {
		return err
		fmt.Printf("handleEvent Error: %+v", err)
	}
	return nil
}

func (d *Discovery) GetContainerByName(name string) (*docker.APIContainers, error) {
	for _, container := range d.containers {
		for _, n := range container.Names {
			if name == n {
				return &container, nil
			}
		}
	}
	return nil, errors.New("Container not found!")
}

func (d *Discovery) GetPortMappings(name string) ([]docker.APIPort, error) {
	cont, err := d.GetContainerByName(name)
	if err != nil {
		return nil, err
	}
	return cont.Ports, nil
}

func (d *Discovery) GetPortMapping(name string, port int64) (int64, error) {
	mappings, err := d.GetPortMappings(name)
	if err != nil {
		return 0, err
	}
	for _, mapping := range mappings {
		if mapping.PrivatePort == port {
			return mapping.PrivatePort, nil
		}
	}
	return 0, errors.New("Port mapping not found!")
}
