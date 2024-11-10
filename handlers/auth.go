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

var sessions = map[string]session{}

// each session contains the username of the user and the time at which it expires
type session struct {
	user   db.User
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
	fmt.Print("[Login] Username and password: ")
	err := r.ParseForm()
	if err != nil {
		panic(err)
		return
	}
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	fmt.Println(username, password)
	user, err := h.svc.GetUserByName(context.Background(), username)
	if err != nil {
		w.Write([]byte("User not found"))
		return
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
		w.Write([]byte("Wrong password"))
		return
	}

	// Create a new random session token
	// we use the "github.com/google/uuid" library to generate UUIDs
	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(120 * time.Second)

	// Set the token in the session map, along with the session information
	sessions[sessionToken] = session{
		user:   *user,
		expiry: expiresAt,
	}

	// Finally, we set the client cookie for "session_token" as the session token we just generated
	// we also set an expiry time of 120 seconds
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: expiresAt,
	})
	http.Redirect(w, r, "/", http.StatusFound)
	fmt.Println("Login successful")
}

func (h *AuthHandler) SignupPost(w http.ResponseWriter, r *http.Request) {
	fmt.Print("[Signup] Username and password: ")
	err := r.ParseForm()
	if err != nil {
		panic(err)
		return
	}
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	fmt.Println(username, password)
	var user = &db.User{
		Name:     username,
		Password: password,
	}
	err = h.svc.CreateUser(context.Background(), user)
	if err != nil {
		panic(err)
		return
	}

	// Create a new random session token
	// we use the "github.com/google/uuid" library to generate UUIDs
	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(120 * time.Second)

	// Set the token in the session map, along with the session information
	sessions[sessionToken] = session{
		user:   *user,
		expiry: expiresAt,
	}

	// Finally, we set the client cookie for "session_token" as the session token we just generated
	// we also set an expiry time of 120 seconds
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: expiresAt,
	})
	http.Redirect(w, r, "/", http.StatusFound)
	fmt.Println("Signup successful")
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

func Signup(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseFiles("signup.html"))
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

func Logout(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")
	if err != nil {
		panic(err)
		return
	}
	delete(sessions, c.Value)
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now(),
	})
	http.Redirect(w, r, "/", http.StatusFound)
}

func MyProfile(w http.ResponseWriter, r *http.Request) {
	// If the session is valid, return the welcome message to the user
	user, err := GetCurrentUser(w, r)
	parameters := map[string]string{}
	if user != nil {
		parameters["current_user"] = user.Name
	} else if err == nil {
		parameters["current_user"] = ""
	} else {
		return
	}
	_, err = w.Write([]byte("Welcome " + user.Name))
	if err != nil {
		return
	}
}

func GetCurrentUser(w http.ResponseWriter, r *http.Request) (*db.User, error) {
	// We can obtain the session token from the requests cookies, which come with every request
	c, err := r.Cookie("session_token")
	if err != nil {
		return nil, err
	}
	sessionToken := c.Value

	// We then get the session from our session map
	userSession, exists := sessions[sessionToken]
	if !exists {
		return nil, err
	}
	if userSession.isExpired() {
		delete(sessions, sessionToken)
		return nil, err
	}
	return &userSession.user, nil
}
