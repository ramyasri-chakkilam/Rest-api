package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"
)

type Todo struct {
	Work string `json:"work"`
	ID   string `json:"id"`
}

var todos = []Todo{
	{Work: "clean room", ID: "ramya"},
	{Work: "sing a song", ID: "mahira"},
}

var mu sync.Mutex

func getTodos(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func createTodo(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	mu.Lock()
	todos = append(todos, todo)
	mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}

func updateTodo(w http.ResponseWriter, r *http.Request, id string) {
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	var updatedTodo Todo
	if err := json.NewDecoder(r.Body).Decode(&updatedTodo); err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for i, todo := range todos {
		if todo.ID == id {
			updatedTodo.ID = id
			todos[i] = updatedTodo

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedTodo)
			return
		}
	}

	http.Error(w, "todo not found", http.StatusNotFound)
}

func deleteTodo(w http.ResponseWriter, r *http.Request, id string) {
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for i, todo := range todos {
		if todo.ID == id {
			todos = append(todos[:i], todos[i+1:]...)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "todo deleted successfully",
			})
			return
		}
	}

	http.Error(w, "todo not found", http.StatusNotFound)
}

func todosHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		getTodos(w, r)
	case "POST":
		createTodo(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func todoByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/todos/")

	switch r.Method {
	case "PUT":
		updateTodo(w, r, id)
	case "DELETE":
		deleteTodo(w, r, id)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/todos", todosHandler)
	http.HandleFunc("/todos/", todoByIDHandler)
	http.ListenAndServe(":8080", nil)
}
