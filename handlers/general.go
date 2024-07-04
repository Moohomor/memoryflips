package handlers

import (
	"net/http"
	// "github.com/go-chi/chi"
	// "context"
	"html/template"
)

func Index(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseFiles("home.html"))
	tpl.Execute(w,nil)
}

func Favicon(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "files/favicon.ico", 302)
}

func List(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseFiles("feed.html"))
	tpl.Execute(w,nil)
}