package handlers

import (
	"context"
	"fmt"
	"html/template"
	"memoryflips/db"
	"net/http"
	"time"

	"github.com/google/uuid"
)

var users = map[int]string{
	1: "password1",
	2: "password2",
}

var sessions = map[string]session{}

// each session contains the username of the user and the time at which it expires
type session struct {
	userId int
	expiry time.Time
}

// we'll use this method later to determine if the session has expired
func (s session) isExpired() bool {
	return s.expiry.Before(time.Now())
}

type Credentials struct {
	Password string `json:"password"`
	UserId   int    `json:"username"`
}

type AuthHandler struct {
	svc *db.Service
}

func NewAuthHandler(svc *db.Service) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) LoginPost(w http.ResponseWriter, r *http.Request) {
	fmt.Print("Username and password: ")
	err := r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	fmt.Println(username, password)
	user, err := h.svc.GetUserByName(context.Background(), username)
	if err != nil {
		panic(err)
	}
	var creds = Credentials{
		UserId:   user.Id,
		Password: password,
	}
	// Get the JSON body and decode into credentials

	// Get the expected password from our in memory map
	expectedPassword := user.Password

	if expectedPassword != creds.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Create a new random session token
	// we use the "github.com/google/uuid" library to generate UUIDs
	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(120 * time.Second)

	// Set the token in the session map, along with the session information
	sessions[sessionToken] = session{
		userId: creds.UserId,
		expiry: expiresAt,
	}

	// Finally, we set the client cookie for "session_token" as the session token we just generated
	// we also set an expiry time of 120 seconds
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: expiresAt,
	})
	http.Redirect(w, r, "", http.StatusFound)
	fmt.Println("Login successful")
}

func Login(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseFiles("login.html"))
	type dt struct {
		message string
	}
	data := dt{
		message: "",
	}
	err := tpl.Execute(w, data)
	if err != nil {
		return
	}
}

func MyProfile(w http.ResponseWriter, r *http.Request) {
	// We can obtain the session token from the requests cookies, which come with every request
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sessionToken := c.Value

	// We then get the session from our session map
	userSession, exists := sessions[sessionToken]
	if !exists {
		// If the session token is not present in session map, return an unauthorized error
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// If the session is present, but has expired, we can delete the session, and return
	// an unauthorized status
	if userSession.isExpired() {
		delete(sessions, sessionToken)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// If the session is valid, return the welcome message to the user
	_, err = w.Write([]byte(fmt.Sprintf("Welcome %i!", userSession.userId)))
	if err != nil {
		return
	}
}
