package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

func (d *Discovery) ViewAPIContainers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	containers, err := json.Marshal(d.containers)
	if err != nil {
		d.log.Printf("JSON marshal error: %v", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(containers)
	return
}

func (d *Discovery) ViewAPIContainersFull(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	contArray := make([]*ProjectContainer, 0, len(d.containersFull))
	for _, c := range d.containersFull {
		contArray = append(contArray, c)
	}
	containers, err := json.Marshal(contArray)
	if err != nil {
		d.log.Printf("JSON marshal error: %v", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(containers)
	return
}

func (d *Discovery) ViewAPIContainerName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	name := ps.ByName("name")
	if len(name) == 0 {
		http.NotFound(w, r)
		return
	}
	container, ok := d.containersFull[name]
	if !ok {
		http.NotFound(w, r)
		return
	}

	containerJSON, err := json.Marshal(container)
	if err != nil {
		d.log.Printf("JSON marshal error: %v", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(containerJSON)
	return

}

func (d *Discovery) ViewAPIContainerMappings(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	name := ps.ByName("name")
	if len(name) == 0 {
		http.NotFound(w, r)
		return
	}

	mappings, err := d.GetPortMappings(name)
	if err != nil {
		d.log.Printf("Could not get port mappings: %v", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	mappingsJSON, err := json.Marshal(mappings)
	if err != nil {
		d.log.Printf("JSON marshal error: %v", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(mappingsJSON)
	return

}

func (d *Discovery) ViewAPIContainerMapping(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	name := ps.ByName("name")
	port := ps.ByName("port")
	protocol := ps.ByName("protocol")
	if len(name) == 0 || len(port) == 0 {
		http.NotFound(w, r)
		return
	}
	if len(protocol) == 0 {
		protocol = "tcp"
	}
	portInt, err := strconv.ParseInt(port, 10, 64)
	if err != nil {
		d.log.Printf("Could not parse port as int: %v", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	iPort := Port{
		Port:     portInt,
		Protocol: protocol,
	}

	mapping, err := d.GetPortMapping(name, iPort)
	if err != nil {
		d.log.Printf("Could not get port mappings: %v", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	getParams := r.URL.Query()
	_, textVersion := getParams["text"]

	if textVersion {
		fmt.Fprintf(w, "%s:%d", d.settings.Hostname, mapping.Port)
		return
	}

	mappingJSON, err := json.Marshal(Mapping{
		Container: iPort,
		Host:      mapping,
		Hostname:  d.settings.Hostname,
	})
	if err != nil {
		d.log.Printf("JSON marshal error: %v", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(mappingJSON)
	return

}

func (d *Discovery) ViewAPIContainerLogs(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	name := ps.ByName("name")
	if len(name) == 0 {
		http.NotFound(w, r)
		return
	}
	container, ok := d.containersFull[name]
	if !ok {
		http.NotFound(w, r)
		return
	}

	logs, err := d.getLogsOfContainer(container, 50000)
	if err != nil {
		d.log.Printf("Could not get logs of container %s: %v", container.FullName, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "%s", logs)

	return

}

func (d *Discovery) ViewAPIProjectUp(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	project := ps.ByName("project")
	if len(project) == 0 {
		http.NotFound(w, r)
		return
	}

	allUp := true
	number := 0

	for _, v := range d.containersFull {
		if v.Project != project {
			continue
		}
		if v.State.Running || v.State.ExitCode == 0 {
			number += 1
			continue
		}
		allUp = false
	}

	if number == 0 {
		http.NotFound(w, r)
		return
	}

	fmt.Fprintf(w, "%t", allUp)

}
