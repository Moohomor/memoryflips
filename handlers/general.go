package handlers

import (
	"errors"
	"net/http"

	// "github.com/go-chi/chi"
	// "context"
	"html/template"
)

func Index(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseFiles("home.html"))
	user, err := GetCurrentUser(w, r)
	parameters := map[string]string{}
	if user != nil {
		parameters["current_user"] = user.Name
	} else if errors.Is(err, http.ErrNoCookie) {
		parameters["current_user"] = ""
	} else if err != nil {
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
	} else if errors.Is(err, http.ErrNoCookie) {
		parameters = map[string]string{
			"current_user": "",
		}
	}

	tpl := template.Must(template.ParseFiles("feed.html"))
	err = tpl.Execute(w, parameters)
	if err != nil {
		return
	}
}
