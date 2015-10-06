package main

import (
	"errors"
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"github.com/julienschmidt/httprouter"
	"github.com/tbruyelle/hipchat-go/hipchat"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

type Discovery struct {
	sync.Mutex
	dockerAPI        string
	listener         chan *docker.APIEvents
	containers       []docker.APIContainers
	containersFull   map[string]*ProjectContainer
	client           *docker.Client
	apiPort          int64
	settings         Settings
	hipchatClient    *hipchat.Client
	log              *log.Logger
	logFile          string
	containerLogPath string
}

func NewDiscovery(dockerAPI string, apiPort int64) (*Discovery, error) {
	d := new(Discovery)
	d.dockerAPI = dockerAPI
	d.apiPort = apiPort
	d.settings = Settings{
		Hostname:         "192.168.178.27",
		Notification:     false,
		SaveLogs:         true,
		SaveLogsDays:     30,
		SaveLogsInterval: 30,
	}
	d.containerLogPath = "logs"
	d.logFile = "discovery.log"
	d.listener = make(chan *docker.APIEvents)
	d.containers = make([]docker.APIContainers, 0)
	d.containersFull = make(map[string]*ProjectContainer)

	lf, err := os.OpenFile(d.logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not open log file %s: %v", d.logFile, err.Error()))
	}
	d.log = log.New(lf, "", log.Ldate|log.Ltime|log.Lshortfile)

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
	d.containersFull = make(map[string]*ProjectContainer)
	for _, container := range containers {
		id := container.ID
		full, err := d.client.InspectContainer(id)
		if err != nil {
			return err
		}
		fullP := NewProjectContainerFromContainer(full)
		d.containersFull[fullP.FullName] = fullP
	}
	return nil
}

func (d *Discovery) handleEvent(event *docker.APIEvents) error {
	d.log.Print(EventToString(event))
	if event.Status == "create" || event.Status == "destroy" || event.Status == "start" || event.Status == "stop" {
		err := d.refreshList()
		if err != nil {
			d.log.Printf("handleEvent Error: %+v", err)
			return err
		}
	}
	if event.Status == "die" && d.settings.Notification {
		err := d.handleCrashEvent(event)
		if err != nil {
			d.log.Printf("handleEvent Error: %+v", err)
			return err
		}
	}
	if event.Status == "destroy" && d.settings.DeleteLogsOnRemove {
		container, err := d.getContainerById(event.ID)
		if err != nil {
			return err
		}
		err = d.deleteLogsOfContainer(container)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *Discovery) handleCrashEvent(event *docker.APIEvents) error {
	container, err := d.getContainerById(event.ID)
	if err != nil {
		return err
	}

	if container.State.ExitCode == 0 {
		return nil
	}

	if d.hipchatClient == nil {
		err := d.initHipChatClient()
		if err != nil {
			return err
		}
	}

	err = d.sendCrashNotification(container)
	if err != nil {
		return err
	}

	return nil
}

func (d *Discovery) gatherLogs() {
	err := os.MkdirAll(d.containerLogPath, 0777)
	if err != nil {
		d.log.Printf("Could not create log folder: %s", err.Error())
		return
	}
	lastLogCheckMap := make(map[string]time.Time)
	for _ = range time.Tick(time.Duration(d.settings.SaveLogsInterval) * time.Second) {
		for _, container := range d.containersFull {
			lastLogCheck, ok := lastLogCheckMap[container.ID]
			if !ok {
				d.gatherLogForContainer(container, 0)
			} else {
				d.gatherLogForContainer(container, lastLogCheck.Unix())
			}
			lastLogCheckMap[container.ID] = time.Now()
		}
	}
}

func (d *Discovery) getContainerById(id string) (*ProjectContainer, error) {
	for _, v := range d.containersFull {
		if v.ID == id {
			return v, nil
		}
	}
	return nil, errors.New("Container not found.")
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
			Hostname: d.settings.Hostname,
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
	r.GET("/api/container/:name/logs", d.ViewAPIContainerLogs)
	r.GET("/api/container/:name/mappings", d.ViewAPIContainerMappings)
	r.GET("/api/container/:name/mapping/:port", d.ViewAPIContainerMapping)
	r.GET("/api/container/:name/mapping/:port/:protocol", d.ViewAPIContainerMapping)
	r.GET("/api/projectUp/:project", d.ViewAPIProjectUp)
	r.GET("/web/settings", d.ViewWebSettings)
	r.POST("/web/settings", d.ViewWebSettings)
	r.GET("/", d.ViewWebContainers)
	r.GET("/web", d.ViewWebContainers)
	r.GET("/web/containers", d.ViewWebContainers)
	r.GET("/web/container/:name/logs", d.ViewWebContainerLogs)
	r.GET("/web/logs", d.ViewWebLogs)
	r.GET("/web/about", d.ViewWebAbout)

	http.ListenAndServe(fmt.Sprintf(":%d", d.apiPort), r)
}
