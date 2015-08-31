package main

import (
	"github.com/julienschmidt/httprouter"
	"html/template"
	"log"
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		hostname := r.PostFormValue("Hostname")
		if len(hostname) == 0 {
			errors = append(errors, "Empty hostname.")
		} else {
			d.settings.Hostname = hostname
		}

		notificationsRaw := r.PostFormValue("Notification")
		if len(notificationsRaw) == 0 {
			d.settings.Notification = false
		} else {
			notifications, err := strconv.ParseBool(notificationsRaw)
			if err != nil {
				log.Print(err.Error())
				errors = append(errors, "Invalid value for notifications.")
			} else {
				d.settings.Notification = notifications
			}
		}

		notificationLogRaw := r.PostFormValue("NotificationLog")
		if len(notificationLogRaw) == 0 {
			d.settings.NotificationLog = false
		} else {
			notificationLog, err := strconv.ParseBool(notificationLogRaw)
			if err != nil {
				log.Print(err.Error())
				errors = append(errors, "Invalid value for \"Add logs to notification\".")
			} else {
				d.settings.NotificationLog = notificationLog
			}
		}

		hipChatToken := r.PostFormValue("HipChatToken")
		if len(hipChatToken) == 0 && d.settings.Notification {
			errors = append(errors, "Empty HipChat API token.")
		} else {
			d.settings.HipChatToken = hipChatToken
		}

		hipChatRoom := r.PostFormValue("HipChatRoom")
		if len(hipChatRoom) == 0 && d.settings.Notification {
			errors = append(errors, "Empty HipChat room.")
		} else {
			d.settings.HipChatRoom = hipChatRoom
		}

		saveLogsRaw := r.PostFormValue("SaveLogs")
		if len(saveLogsRaw) == 0 {
			d.settings.SaveLogs = false
		} else {
			saveLogs, err := strconv.ParseBool(saveLogsRaw)
			if err != nil {
				log.Print(err.Error())
				errors = append(errors, "Invalid value for save logs.")
			} else {
				d.settings.SaveLogs = saveLogs
			}
		}

		saveLogsDays, err := strconv.ParseInt(r.PostFormValue("SaveLogsDays"), 10, 64)
		if err != nil {
			log.Print(err.Error())
			errors = append(errors, "Invalid value for save log days.")
		}
		d.settings.SaveLogsDays = int(saveLogsDays)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tpl.Execute(w, tplData); err != nil {
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tpl.Execute(w, tplData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	return

}
