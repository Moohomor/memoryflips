package handlers

import (
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
	user_id int
	expiry  time.Time
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
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	fmt.Println(username, password)
	h.svc.GetUserByName(username)
	var creds Credentials = Credentials{
		UserId:   1,
		Password: password,
	}
	// Get the JSON body and decode into credentials

	// Get the expected password from our in memory map
	expectedPassword, ok := users[creds.UserId]

	// If a password exists for the given user
	// AND, if it is the same as the password we received, the we can move ahead
	// if NOT, then we return an "Unauthorized" status
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
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
		user_id: creds.UserId,
		expiry:  expiresAt,
	}

	// Finally, we set the client cookie for "session_token" as the session token we just generated
	// we also set an expiry time of 120 seconds
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: expiresAt,
	})
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
	tpl.Execute(w, data)
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
	w.Write([]byte(fmt.Sprintf("Welcome %i!", userSession.user_id)))
}
