package handlers

import (
	"errors"
	"memoryflips/db"
	"net/http"
	"strconv"
	// "github.com/go-chi/chi"
	"context"
	"html/template"
)

type WordsHandler struct {
	svc *db.Service
}

func NewWordsHandler(svc *db.Service) *WordsHandler {
	return &WordsHandler{svc: svc}
}

func (h *WordsHandler) Question(w http.ResponseWriter, r *http.Request) {
	word, err := h.svc.GetWord(context.Background())
	wordCount, err := GetWordCounter(w, r)
	if wordCount >= 15 {
		http.Redirect(w, r, "/result", http.StatusFound)
	}
	if wordCount == 0 {
		http.SetCookie(w, &http.Cookie{
			Name:  "word_counter",
			Value: "0",
		})
	}
	if err != nil {
		panic(err)
		return
	}
	var parameters map[string]string
	if word != nil {
		parameters = map[string]string{
			"word_id": strconv.Itoa(word.Id),
			"word":    word.Rus,
		}
	}
	if err != nil {
		panic(err)
		return
	}
	wordCount++
	tpl := template.Must(template.ParseFiles("question.html"))
	err = tpl.Execute(w, parameters)
	if err != nil {
		return
	}
}

func (h *WordsHandler) Answer(w http.ResponseWriter, r *http.Request) {
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
	tpl := template.Must(template.ParseFiles("answer.html"))
	err = tpl.Execute(w, parameters)
	if err != nil {
		return
	}
}

func GetWordCounter(w http.ResponseWriter, r *http.Request) (int, error) {
	// We can obtain the session token from the requests cookies, which come with every request
	c, err := r.Cookie("word_counter")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			http.SetCookie(w, &http.Cookie{
				Name:  "word_counter",
				Value: "0",
			})
		}
		return 0, err
	}
	i, err := strconv.Atoi(c.Value)
	if err != nil {
		panic(err)
		return 0, nil
	}
	return i, nil
}
