package main

import (
	"errors"
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type Discovery struct {
	sync.Mutex
	dockerAPI      string
	listener       chan *docker.APIEvents
	containers     []docker.APIContainers
	containersFull map[string]*docker.Container
	client         *docker.Client
	apiPort        int64
	hostname       string
}

func NewDiscovery(dockerAPI string, apiPort int64) (*Discovery, error) {
	d := new(Discovery)
	d.dockerAPI = dockerAPI
	d.apiPort = apiPort
	d.hostname = "marb.ec"
	d.listener = make(chan *docker.APIEvents)
	d.containers = make([]docker.APIContainers, 0)
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
	d.Lock()
	defer d.Unlock()
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

func (d *Discovery) GetPortMappings(name string) ([]Mapping, error) {
	container, ok := d.containersFull[name]
	if !ok {
		return nil, errors.New("Container not found!")
	}

	mappings := make([]Mapping, 0, len(container.NetworkSettings.Ports))
	for k, v := range container.NetworkSettings.Ports {
		iPort, _ := strconv.ParseInt(k.Port(), 10, 64)
		if len(v) == 0 {
			continue
		}
		host := v[0]
		hPort, _ := strconv.ParseInt(host.HostPort, 10, 64)
		mappings = append(mappings, Mapping{
			Container: Port{
				Port:     iPort,
				Protocol: k.Proto(),
			},
			Host: Port{
				Port:     hPort,
				Protocol: k.Proto(),
			},
			Hostname: d.hostname,
		})
	}
	return mappings, nil
}

func (d *Discovery) GetPortMapping(name string, port Port) (Port, error) {
	mappings, err := d.GetPortMappings(name)
	if err != nil {
		return Port{}, err
	}

	for _, mapping := range mappings {
		if mapping.Container.String() == port.String() {
			return mapping.Host, nil
		}
	}

	return Port{}, errors.New(fmt.Sprintf("No port mapping available for internal port %s", port))

}

func (d *Discovery) serveWeb() {
	r := httprouter.New()
	r.GET("/api/containers", d.ViewAPIContainers)
	r.GET("/api/containersFull", d.ViewAPIContainersFull)
	r.GET("/api/container/:name", d.ViewAPIContainerName)
	r.GET("/api/container/:name/mappings", d.ViewAPIContainerMappings)
	r.GET("/api/container/:name/mapping/:port", d.ViewAPIContainerMapping)
	r.GET("/api/container/:name/mapping/:port/:protocol", d.ViewAPIContainerMapping)
	r.GET("/web/settings", d.ViewWebSettings)
	r.POST("/web/settings", d.ViewWebSettingsPost)
	r.GET("/web/containers", d.ViewWebContainers)

	http.ListenAndServe(fmt.Sprintf(":%d", d.apiPort), r)
}