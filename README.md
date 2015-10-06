[![Build Status](https://travis-ci.org/herzog31/service-discovery.svg?branch=master)](https://travis-ci.org/herzog31/service-discovery)
[![GoDoc](https://godoc.org/github.com/herzog31/service-discovery?status.svg)](https://godoc.org/github.com/herzog31/service-discovery)
[![Docker Hub](https://img.shields.io/docker/pulls/herzog31/service-discovery.svg)](https://hub.docker.com/r/herzog31/service-discovery)
[![Release](https://img.shields.io/github/release/herzog31/service-discovery.svg)](https://github.com/herzog31/service-discovery/releases)
[![Go](https://img.shields.io/badge/Go-1.5.1-blue.svg)](https://golang.org/)

# service-discovery
Service Discovery for Docker written in Go

## Usage
Using the API of your Docker Engine the service discovery allows you to access valuable information.
These information are available either via the web interface at `localhost:8080/web/containers` or the REST API.
In the web interface you can see a list of all containers and their logs, a settings page and the logs of the service discovery.

### Settings
First, please set the hostname of Docker host, so the service discovery can generate valid links to your containers.

You can enable notifications, so that whenever one of your containers crashes, you receive a notification via [HipChat](https://www.hipchat.com/).
It is possible to include recent logging information of the crashed container in the notification.

In the last section of the settings, you can configure the logging behaviour of the service discovery.

### Requirements
Requirement     | Version
--------------- | -------------
Docker          | 1.8.1+
Docker-Compose  | 1.4.0+

### Native
Download an executable that matches your architecture and the `template` folder from the repository.
Make the file executable. Example for Linux:
```
chmod +x linux_amd64_service-discovery
```

You can configure the port of the web interface and the REST API and the location of the Docker API socket via flags.
For more information on the flags, execute `service-discovery --help`.

### Docker Container
The service discovery is also available as Docker container for Linux (amd64) from [Docker Hub](https://hub.docker.com/r/herzog31/service-discovery/).
You can start it via the following command:
```
docker run -d -p 8080:8080 -v /var/run/docker.sock:/var/run/docker.sock --name discovery herzog31/service-discovery:latest
```

In the repository you also find a Dockerfile to build a container for Linux (ARM).

## API Documentation
Please refer to [API.md](API.md) for a full API documentation.