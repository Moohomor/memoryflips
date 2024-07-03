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

	"hse_school/internal/handlers"
	// "hse_school/internal/repositories"
	// "hse_school/internal/services"
)

func main() {
	// Подключение к базе данных
	db, err := pgxpool.Connect(context.Background(), "postgresql://postgres:password@127.0.0.1/postgres")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Создание репозитория, сервиса и обработчиков
	// repo := repositories.NewUserRepository(db)
	// service := services.NewDefaultUserService(repo)
	// handler := handlers.NewUserHandler(service)

	// Настройка маршрутизатора
	r := chi.NewRouter()

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "cmd/project/data"))
	FileServer(r, "/files", filesDir)

	r.Get("/", handlers.Index)
	r.Get("/favicon.ico", handlers.Favicon)
	// r.Get("/users/{id}", handler.GetUserByID)
	// r.Post("/users", handler.CreateUser)
	
	// Запуск HTTP сервера
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}


func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
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