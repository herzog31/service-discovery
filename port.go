package main

import (
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"strconv"
)

type Mapping struct {
	Container Port
	Host      Port
	Hostname  string
}

func (m Mapping) String() string {
	if m.Host.Protocol != "tcp" {
		return fmt.Sprintf("%s:%d/%s", m.Hostname, m.Host.Port, m.Host.Protocol)
	} else {
		return fmt.Sprintf("%s:%d", m.Hostname, m.Host.Port)
	}
}

type Port struct {
	Port     int64
	Protocol string
}

func NewPortFromDockerPort(dPort docker.Port) Port {
	parsed, _ := strconv.ParseInt(dPort.Port(), 10, 64)
	return Port{
		Port:     parsed,
		Protocol: dPort.Proto(),
	}
}

func (p Port) ToDockerPort() docker.Port {
	return docker.Port(fmt.Sprintf("%d/%s", p.Port, p.Protocol))
}

func (p Port) String() string {
	return fmt.Sprintf("%d/%s", p.Port, p.Protocol)
}
