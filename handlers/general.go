package handlers

import (
	"net/http"
	// "github.com/go-chi/chi"
	// "context"
	"html/template"
)

func Index(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseFiles("home.html"))
	err := tpl.Execute(w, nil)
	if err != nil {
		return
	}
}

func Favicon(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "files/favicon.ico", http.StatusFound)
}

func List(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseFiles("feed.html"))
	err := tpl.Execute(w, nil)
	if err != nil {
		return
	}
}
