package handlers

import (
	"net/http"

	// "github.com/go-chi/chi"
	// "context"
	"html/template"
)

func Index(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseFiles("home.html"))
	user, err := GetCurrentUser(w, r)
	var parameters map[string]string
	if user != nil {
		parameters = map[string]string{
			"current_user": user.Name,
		}
	} else if err != nil {
		parameters = map[string]string{
			"current_user": "",
		}
	} else {
		return
	}
	err = tpl.Execute(w, parameters)
	if err != nil {
		return
	}
}

func Favicon(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "files/favicon.ico", http.StatusFound)
}

func Feed(w http.ResponseWriter, r *http.Request) {
	user, err := GetCurrentUser(w, r)
	var parameters map[string]string
	if user != nil {
		parameters = map[string]string{
			"current_user": user.Name,
		}
	} else if err != nil {
		parameters = map[string]string{
			"current_user": "",
		}
	} else {
		return
	}
	tpl := template.Must(template.ParseFiles("feed.html"))
	err = tpl.Execute(w, parameters)
	if err != nil {
		return
	}
}
