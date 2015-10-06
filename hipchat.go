package main

import (
	"errors"
	"fmt"
	"github.com/tbruyelle/hipchat-go/hipchat"
	"strings"
)

func (d *Discovery) initHipChatClient() error {

	if len(d.settings.HipChatToken) == 0 {
		return errors.New("No HipChat token specified.")
	}

	d.hipchatClient = hipchat.NewClient(d.settings.HipChatToken)
	hipchat.AuthTest = true
	d.hipchatClient.Room.List()
	_, ok := hipchat.AuthTestResponse["success"]
	hipchat.AuthTest = false

	if !ok {
		return errors.New("Invalid HipChat token.")
	}

	return nil
}

func (d *Discovery) sendCrashNotification(container *ProjectContainer) error {

	if len(d.settings.HipChatRoom) == 0 {
		return errors.New("No HipChat room specified.")
	}

	rooms, _, err := d.hipchatClient.Room.List()
	if err != nil {
		return err
	}

	var message string
	if container.Project != "" {
		message = fmt.Sprintf("Container <b>%s</b> of project <b>%s</b> crashed at %v with exit code <b>%d</b>.", container.FullName, container.Project, container.State.FinishedAt.Local(), container.State.ExitCode)
	} else {
		message = fmt.Sprintf("Container <b>%s</b> crashed at %v with exit code <b>%d</b>.", container.FullName, container.State.FinishedAt.Local(), container.State.ExitCode)
	}

	if d.settings.NotificationLog {
		logs, err := d.getLogsOfContainer(container, 5000)
		message += "<br />"
		if err == nil {
			logs = strings.Replace(logs, "\n", "<br />", -1)
			message += logs
		} else {
			d.log.Printf("Could not get logs of container %s: %v", container.FullName, err.Error())
		}
	}

	notification := &hipchat.NotificationRequest{
		Color:         "red",
		Notify:        true,
		MessageFormat: "html",
		Message:       message,
	}

	for _, room := range rooms.Items {
		if room.Name == d.settings.HipChatRoom {
			_, err := d.hipchatClient.Room.Notification(room.Name, notification)
			if err != nil {
				return err
			}
			return nil
		}
	}

	return errors.New("HipChat room not found!")
}
