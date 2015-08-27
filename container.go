package main

import (
	"github.com/fsouza/go-dockerclient"
	"strconv"
	"strings"
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

	/* container.Container = *cont
	if strings.Contains(container.Name, "_") {
		container.Project = strings.SplitN(container.Name, "_", 2)[0]
		container.Name = strings.SplitN(container.Name, "_", 2)[1]
	} else {
		container.Project = "none"
	} */

	return container
}
