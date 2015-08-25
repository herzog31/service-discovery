package main

import (
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"strconv"
)

type Port struct {
	port     int64
	protocol string
}

func NewPortFromDockerPort(dPort docker.Port) Port {
	parsed, _ := strconv.ParseInt(dPort.Port(), 10, 64)
	return Port{
		port:     parsed,
		protocol: dPort.Proto(),
	}
}

func (p Port) ToDockerPort() docker.Port {
	return docker.Port(fmt.Sprintf("%d/%s", p.port, p.protocol))
}

func (p Port) String() string {
	return fmt.Sprintf("%d/%s", p.port, p.protocol)
}
