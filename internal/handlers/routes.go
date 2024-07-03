package handlers

import (
	"net/http"
	// "github.com/go-chi/chi"
	// "context"
	"html/template"
)

func Index(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseFiles("cmd/project/home.html"))
	tpl.Execute(w,nil)
}

func Favicon(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "files/favicon.ico", 302)
}