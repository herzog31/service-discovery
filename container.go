package main

import (
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"strconv"
	"strings"
	"time"
)

type ProjectContainer struct {
	docker.Container
	FullName string
	Project  string
	Number   int64
}

func NewProjectContainerFromContainer(cont *docker.Container) *ProjectContainer {
	container := new(ProjectContainer)
	container.Container = *cont
	container.Number = 1
	container.Project = ""
	container.FullName = strings.TrimPrefix(container.Name, "/")

	project, okp := container.Config.Labels["com.docker.compose.project"]
	name, okn := container.Config.Labels["com.docker.compose.service"]
	if okp && okn {
		container.Project = project
		container.Name = name
	} else {
		container.Name = container.FullName
	}

	if number, ok := container.Config.Labels["com.docker.compose.container-number"]; ok {
		if pNumber, err := strconv.ParseInt(number, 10, 64); err == nil {
			container.Number = pNumber
		}
	}

	return container
}

func (p *ProjectContainer) TplGetCommand() string {
	return strings.Join(p.Config.Cmd, " ")
}

type TplPort struct {
	Exposed string
	Mapping string
}

func EventToString(e *docker.APIEvents) string {
	return fmt.Sprintf("Event: %s of container with ID %s (Image: %s) at %v", e.Status, e.ID, e.From, time.Unix(e.Time, 0).Local())
}

func (p *ProjectContainer) TplGetPorts() []TplPort {

	ports := make([]TplPort, 0, len(p.NetworkSettings.Ports))
	for k, v := range p.NetworkSettings.Ports {
		if len(v) == 0 {
			ports = append(ports, TplPort{Exposed: string(k)})
			continue
		}
		host := v[0]
		ports = append(ports, TplPort{
			Exposed: string(k),
			Mapping: host.HostPort,
		})
	}

	return ports
}
