package main

import (
	"github.com/julienschmidt/httprouter"
	"html/template"
	"net/http"
	"path"
	"strconv"
	"strings"
)

func (d *Discovery) ViewWebSettings(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	errors := make([]string, 0)

	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			d.log.Printf("Could not parse form: %v", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		hostname := r.PostFormValue("Hostname")
		if len(hostname) == 0 {
			errors = append(errors, "Empty hostname.")
		}
		if d.settings.Hostname != hostname {
			d.settings.Hostname = hostname
			d.log.Printf("Settings: Hostname updated to %s", hostname)
		}

		notificationsRaw := r.PostFormValue("Notification")
		if len(notificationsRaw) == 0 {
			notificationsRaw = "false"
		}
		notifications, err := strconv.ParseBool(notificationsRaw)
		if err != nil {
			errors = append(errors, "Invalid value for notifications.")
		}
		if d.settings.Notification != notifications {
			d.settings.Notification = notifications
			d.log.Printf("Settings: Notifications set to %t", notifications)
		}

		notificationLogRaw := r.PostFormValue("NotificationLog")
		if len(notificationLogRaw) == 0 {
			notificationLogRaw = "false"
		}
		notificationLog, err := strconv.ParseBool(notificationLogRaw)
		if err != nil {
			errors = append(errors, "Invalid value for \"Add logs to notification\".")
		}
		if d.settings.NotificationLog != notificationLog {
			d.settings.NotificationLog = notificationLog
			d.log.Printf("Settings: NotificationLog set to %t", notificationLog)
		}

		hipChatToken := r.PostFormValue("HipChatToken")
		if len(hipChatToken) == 0 && d.settings.Notification {
			errors = append(errors, "Empty HipChat API token.")
		} else {
			if d.settings.HipChatToken != hipChatToken {
				d.settings.HipChatToken = hipChatToken
				d.log.Printf("Settings: New HipChat token set")
				err := d.initHipChatClient()
				if err != nil {
					errors = append(errors, "Invalid HipChat API token.")
				}
			}
		}

		hipChatRoom := r.PostFormValue("HipChatRoom")
		if len(hipChatRoom) == 0 && d.settings.Notification {
			errors = append(errors, "Empty HipChat room.")
		}
		if d.settings.HipChatRoom != hipChatRoom {
			d.settings.HipChatRoom = hipChatRoom
			d.log.Printf("Settings: HipChat room set to %s", hipChatRoom)
		}

		saveLogsRaw := r.PostFormValue("SaveLogs")
		if len(saveLogsRaw) == 0 {
			saveLogsRaw = "false"
		}
		saveLogs, err := strconv.ParseBool(saveLogsRaw)
		if err != nil {
			errors = append(errors, "Invalid value for save logs.")
		}
		if d.settings.SaveLogs != saveLogs {
			d.settings.SaveLogs = saveLogs
			d.log.Printf("Settings: SaveLogs set to %t", saveLogs)
		}

		saveLogsDays, err := strconv.ParseInt(r.PostFormValue("SaveLogsDays"), 10, 64)
		if err != nil {
			errors = append(errors, "Invalid value for save log days.")
		}
		if d.settings.SaveLogsDays != saveLogsDays {
			d.settings.SaveLogsDays = saveLogsDays
			d.log.Printf("Settings: SaveLogDays set to %d", d.settings.SaveLogsDays)
		}
	}

	tplData := struct {
		RequestURL string
		Settings   Settings
		Errors     []string
	}{
		r.URL.RequestURI(),
		d.settings,
		errors,
	}

	layoutPath := path.Join("templates", "layout.html")
	settingsPath := path.Join("templates", "settings.html")

	tpl, err := template.ParseFiles(layoutPath, settingsPath)
	if err != nil {
		d.log.Printf("Could not parse template: %v", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tpl.Execute(w, tplData); err != nil {
		d.log.Printf("Could not execute template: %v", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	return

}

func (d *Discovery) ViewWebContainers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	containersPerProject := make(map[string][]*ProjectContainer)
	for _, v := range d.containersFull {
		if len(v.Project) == 0 {
			containersPerProject["Individual"] = append(containersPerProject["Individual"], v)
		} else {
			project := strings.ToUpper(v.Project[:1]) + v.Project[1:]
			containersPerProject[project] = append(containersPerProject[project], v)
		}
	}

	tplData := struct {
		Containers map[string][]*ProjectContainer
		RequestURL string
		Hostname   string
	}{
		containersPerProject,
		r.URL.RequestURI(),
		d.settings.Hostname,
	}

	layoutPath := path.Join("templates", "layout.html")
	containersPath := path.Join("templates", "containers.html")

	tpl, err := template.ParseFiles(layoutPath, containersPath)
	if err != nil {
		d.log.Printf("Could not parse template: %v", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tpl.Execute(w, tplData); err != nil {
		d.log.Printf("Could not execute template: %v", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	return

}

func (d *Discovery) ViewWebLogs(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	logs, err := d.getLogs(50000)
	if err != nil {
		d.log.Printf("Could not get service discovery logs: %v", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tplData := struct {
		RequestURL string
		Hostname   string
		Logs       string
	}{
		r.URL.RequestURI(),
		d.settings.Hostname,
		logs,
	}

	layoutPath := path.Join("templates", "layout.html")
	containersPath := path.Join("templates", "logs.html")

	tpl, err := template.ParseFiles(layoutPath, containersPath)
	if err != nil {
		d.log.Printf("Could not parse template: %v", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tpl.Execute(w, tplData); err != nil {
		d.log.Printf("Could not execute template: %v", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	return

}

func (d *Discovery) ViewWebContainerLogs(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

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

	tplData := struct {
		RequestURL    string
		Hostname      string
		Logs          string
		ContainerName string
	}{
		r.URL.RequestURI(),
		d.settings.Hostname,
		logs,
		container.FullName,
	}

	layoutPath := path.Join("templates", "layout.html")
	containersPath := path.Join("templates", "containerLogs.html")

	tpl, err := template.ParseFiles(layoutPath, containersPath)
	if err != nil {
		d.log.Printf("Could not parse template: %v", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tpl.Execute(w, tplData); err != nil {
		d.log.Printf("Could not execute template: %v", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	return

}
