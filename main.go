package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"time"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

var DB *pgxpool.Pool

func getRoot(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, nil)
}

type Todo struct {
	Id         uuid.UUID
	Content    string
	Finished   bool
	Created_at time.Time
	Updated_at time.Time
}

func getTodos(w http.ResponseWriter, r *http.Request) {
	var todos []*Todo
	pgxscan.Select(context.Background(), DB, &todos, `SELECT * FROM todos ORDER BY created_at`)
	tmpl := template.Must(template.ParseFiles("templates/todo.html"))
	tmpl.Execute(w, todos)
}

func toggleTodo(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	now := time.Now().UTC()

	res, err := DB.Exec(context.Background(), `
    UPDATE todos
    SET finished = NOT finished, updated_at = $1
    WHERE id = $2
    `, now, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if res.RowsAffected() == 0 {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	var todo []*Todo
	pgxscan.Select(context.Background(), DB, &todo, `SELECT * FROM todos WHERE id = $1`, id)

	tmpl := template.Must(template.ParseFiles("templates/todo.html"))
	tmpl.Execute(w, todo)
}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	res, err := DB.Exec(context.Background(), `DELETE FROM todos WHERE id = $1`, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if res.RowsAffected() == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

func newTodo(content string) Todo {
	now := time.Now().UTC()
	return Todo{
		Id:         uuid.New(),
		Content:    content,
		Finished:   false,
		Created_at: now,
		Updated_at: now,
	}
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	content := r.FormValue("content")
	if content == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	todo := newTodo(content)
	_, err := DB.Exec(context.Background(),
		`INSERT INTO todos (id, content, finished, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)`,
		todo.Id, todo.Content, todo.Finished, todo.Created_at, todo.Updated_at)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	todos := []Todo{todo}

	tmpl := template.Must(template.ParseFiles("templates/todo.html"))
	tmpl.Execute(w, todos)
}

func main() {
	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		fmt.Fprintf(os.Stderr, "No DATABASE_URL env var!")
		os.Exit(1)
	}

	var err error
	DB, err = pgxpool.Connect(context.Background(), dbUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not connect to DB")
		os.Exit(1)
	}
	defer DB.Close()

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Get("/", getRoot)
	router.Get("/todos", getTodos)
	router.Put("/todos/{id}", toggleTodo)
	router.Delete("/todos/{id}", deleteTodo)
	router.Post("/todos", createTodo)
	http.ListenAndServe("0.0.0.0:8080", router)
}
