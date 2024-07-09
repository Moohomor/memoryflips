package handlers

import (
	"errors"
	"fmt"
	"memoryflips/db"
	"net/http"
	"strconv"
	"time"

	// "github.com/go-chi/chi"
	"context"
	"html/template"

	"github.com/go-chi/chi"
)

type WordsHandler struct {
	svc *db.Service
}

func NewWordsHandler(svc *db.Service) *WordsHandler {
	return &WordsHandler{svc: svc}
}

func (h *WordsHandler) Question(w http.ResponseWriter, r *http.Request) {
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

	// total words count
	wordCount, err := GetWordCounter(w, r)
	if wordCount >= 15 {
		http.Redirect(w, r, "/result", http.StatusFound)
		return
	}
	if wordCount == 0 {
		http.SetCookie(w, &http.Cookie{
			Name:  "word_counter",
			Value: "0",
		})
	}
	http.SetCookie(w, &http.Cookie{
		Name:  "word_counter",
		Value: strconv.Itoa(wordCount + 1),
	})
	if err != nil {
		panic(err)
		return
	}

	word, err := h.svc.GetWord(context.Background())
	if err != nil {
		panic(err)
		return
	}
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
	} else if errors.Is(err, http.ErrNoCookie) {
		parameters = map[string]string{
			"current_user": "",
		}
	} else if err != nil {
		panic(err)
		return
	}

	idParam, err := strconv.Atoi(chi.URLParam(r, "wid"))
	if err != nil {
		panic(err)
		return
	}
	word, err := h.svc.GetWordById(context.Background(), idParam)
	if err != nil {
		panic(err)
		return
	}
	if word != nil {
		parameters["word_id"] = strconv.Itoa(word.Id)
		parameters["word"] = word.Eng
	}
	tpl := template.Must(template.ParseFiles("answer.html"))
	err = tpl.Execute(w, parameters)
	if err != nil {
		return
	}
}

func (h *WordsHandler) Result(w http.ResponseWriter, r *http.Request) {
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

	http.SetCookie(w, &http.Cookie{
		Name:  "word_counter",
		Value: "0",
	})

	correctCount, err := GetCorrectCounter(w, r)
	if err != nil {
		panic(err)
		return
	}
	if correctCount == 0 {
		http.SetCookie(w, &http.Cookie{
			Name:  "correct_counter",
			Value: "0",
		})
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "correct_counter",
		Value: "0",
	})

	wordCount, err := GetWordCounter(w, r)
	if wordCount == 0 {
		http.Redirect(w, r, "/feed", http.StatusFound)
		return
	}
	if err != nil {
		panic(err)
		return
	}

	percentage := fmt.Sprintf("%v", int(float32(correctCount)/float32(wordCount)*1000)/10)
	parameters["percentage"] = percentage
	parameters["correct"] = strconv.Itoa(correctCount)
	parameters["total"] = strconv.Itoa(wordCount)

	tpl := template.Must(template.ParseFiles("result.html"))
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
			expiresAt := time.Now().Add(120 * time.Hour)
			http.SetCookie(w, &http.Cookie{
				Name:    "word_counter",
				Value:   "0",
				Expires: expiresAt,
			})
			return 0, nil
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

func GetCorrectCounter(w http.ResponseWriter, r *http.Request) (int, error) {
	// We can obtain the session token from the requests cookies, which come with every request
	c, err := r.Cookie("correct_counter")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			expiresAt := time.Now().Add(120 * time.Hour)
			http.SetCookie(w, &http.Cookie{
				Name:    "correct_counter",
				Value:   "0",
				Expires: expiresAt,
			})
			return 0, nil
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

func (h *WordsHandler) Learned(w http.ResponseWriter, r *http.Request) {
	// correct answers (words) count
	correctCount, err := GetCorrectCounter(w, r)
	if correctCount == 0 {
		http.SetCookie(w, &http.Cookie{
			Name:  "correct_counter",
			Value: "0",
		})
	}
	http.SetCookie(w, &http.Cookie{
		Name:  "correct_counter",
		Value: strconv.Itoa(correctCount + 1),
	})
	if err != nil {
		panic(err)
		return
	}

	// mark word as learned
	wid := chi.URLParam(r, "wid")
	fmt.Println(wid)
	http.Redirect(w, r, "/q", http.StatusFound)
}
