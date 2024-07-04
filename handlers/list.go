package handlers

import (
	"net/http"
	// "github.com/go-chi/chi"
	// "context"
	"html/template"
	"fmt"
)

func Question(w http.ResponseWriter, r *http.Request) {
	fmt.Println("q")
	tpl := template.Must(template.ParseFiles("question.html"))
	tpl.Execute(w,nil)
}