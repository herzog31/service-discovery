package main

import (
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"os"
	"path"
	"strings"
)

func (d *Discovery) gatherLogForContainer(container *ProjectContainer, since int64) {
	logPath := path.Join(d.containerLogPath, fmt.Sprintf("%s.log", container.FullName))
	lf, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		d.log.Printf("Could not open logfile for container %s: %s", container.FullName, err.Error())
		return
	}
	logOptions := docker.LogsOptions{
		Container:    container.ID,
		OutputStream: lf,
		ErrorStream:  lf,
		Stdout:       true,
		Stderr:       true,
		Follow:       false,
		Since:        since,
		Timestamps:   true,
	}
	err = d.client.Logs(logOptions)
	if err != nil {
		d.log.Printf("Could not gather logs of container %s: %s", container.FullName, err.Error())
	}
	defer lf.Close()
	return
}

func (d *Discovery) getLogs(maxSize int64) (string, error) {
	return readLastBytesOfFile(d.logFile, maxSize)
}

func readLastBytesOfFile(path string, maxSize int64) (string, error) {
	lf, err := os.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		return "", err
	}
	defer lf.Close()

	stats, err := lf.Stat()
	start := int64(0)
	bufferSize := stats.Size()
	if stats.Size() > maxSize {
		bufferSize = maxSize
		start = stats.Size() - bufferSize
	}
	logBuffer := make([]byte, bufferSize)

	_, err = lf.ReadAt(logBuffer, start)
	if err != nil {
		return "", err
	}

	logs := string(logBuffer)
	firstBreak := strings.Index(logs, "\n")
	logs = logs[firstBreak:]
	logs = strings.TrimSpace(logs)

	return logs, nil
}

func (d *Discovery) getLogsOfContainer(container *ProjectContainer, maxSize int64) (string, error) {
	return readLastBytesOfFile(path.Join(d.containerLogPath, fmt.Sprintf("%s.log", container.FullName)), maxSize)
}

func (d *Discovery) deleteLogsOfContainer(container *ProjectContainer) error {
	return os.Remove(path.Join(d.containerLogPath, fmt.Sprintf("%s.log", container.FullName)))
}
