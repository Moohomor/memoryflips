package main

import (
	// "fmt"

	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v4/pgxpool"

	dbsvc "memoryflips/db"
	"memoryflips/handlers"
)

func main() {
	// Подключение к базе данных
	db, err := pgxpool.Connect(context.Background(), "postgresql://postgres:password@127.0.0.1/memoryflip")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Настройка маршрутизатора
	r := chi.NewRouter()

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "data"))
	FileServer(r, "/files", filesDir)

	authHandlers := handlers.NewAuthHandler(dbsvc.NewService(db))
	listHandlers := handlers.NewWordsHandler(dbsvc.NewService(db))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		_, err = w.Write([]byte("route does not exist"))
		if err != nil {
			return
		}
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(405)
		_, err = w.Write([]byte("method is not valid"))
		if err != nil {
			return
		}
	})

	r.Get("/", handlers.Index)
	r.Get("/favicon.ico", handlers.Favicon)
	r.Get("/feed", handlers.Feed)

	r.Get("/q", listHandlers.Question)
	//r.Get("/a-{wid}", listHandlers.Answer)
	r.Get("/l-{wid}", listHandlers.Learned)
	r.Get("/result", listHandlers.Result)

	r.Get("/login", handlers.Login)
	r.Post("/login", authHandlers.LoginPost)
	r.Get("/signup", handlers.Signup)
	r.Post("/signup", authHandlers.SignupPost)
	r.Get("/logout", handlers.Logout)
	r.Get("/me", handlers.MyProfile)

	// Запуск HTTP сервера
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe("192.168.1.128:8080", r))
}

func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
