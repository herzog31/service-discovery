package main

import (
	"github.com/julienschmidt/httprouter"
	"html/template"
	"net/http"
	"path"
)

func (d *Discovery) ViewWebSettings(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	tplData := struct {
		RequestURL string
	}{
		r.URL.RequestURI(),
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

func (d *Discovery) ViewWebSettingsPost(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

}

func (d *Discovery) ViewWebContainers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	containersPerProject := make(map[string][]*ProjectContainer)
	for _, v := range d.containersFull {
		if len(v.Project) == 0 {
			containersPerProject["none"] = append(containersPerProject["none"], v)
		} else {
			containersPerProject[v.Project] = append(containersPerProject[v.Project], v)
		}
	}

	tplData := struct {
		Containers map[string][]*ProjectContainer
		RequestURL string
	}{
		containersPerProject,
		r.URL.RequestURI(),
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
