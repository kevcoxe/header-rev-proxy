package handlers

import (
	"fmt"
	"log"
	"net/http"

	"header-rev-proxy/components"

	"github.com/a-h/templ"
)

func onError(w http.ResponseWriter, err error, msg string, code int) {
	if err != nil {
		http.Error(w, msg, code)
		log.Println(msg, err)
	}
}

func RenderView(w http.ResponseWriter, r *http.Request, view templ.Component, layoutPath string) {
	if r.Header.Get("Hx-Request") == "true" {
		err := view.Render(r.Context(), w)
		onError(w, err, "Internal server error", http.StatusInternalServerError)
		return
	}
	err := components.Layout(layoutPath).Render(r.Context(), w)
	onError(w, err, "Internal server error", http.StatusInternalServerError)
}

func HomeGetHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HomeGetHandler")
	RenderView(w, r, components.HomeView("hello, world!"), "/")
}

func TimePostHandler(w http.ResponseWriter, r *http.Request) {
	clientTimeStr := r.FormValue("time")
	err := components.Time(clientTimeStr).Render(r.Context(), w)
	onError(w, err, "Internal server error", http.StatusInternalServerError)
}

func LoginPostHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")

	// create a cookie, set the value, redirect to /grafana route
	// set cookie path to /grafana
	cookie := http.Cookie{
		Name:  "header-rev-proxy-username",
		Value: username,
		Path:  "/grafana",
	}
	http.SetCookie(w, &cookie)

	http.Redirect(w, r, "/grafana/", http.StatusSeeOther)
}
