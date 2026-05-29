package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"
)

type x struct {
	Work string `json:"work"`
	ID   string `json:"id"`
}

var y = []x{
	{Work: "clean room", ID: "ramya"},
	{Work: "sing a song", ID: "mahira"},
}

var mu sync.Mutex

func handlerGet(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(y)
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	var z x
	if err := json.NewDecoder(r.Body).Decode(&z); err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	mu.Lock()
	y = append(y, z)
	mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(z)
}

func putWork(w http.ResponseWriter, r *http.Request, id string) {
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	var updatedWork x
	if err := json.NewDecoder(r.Body).Decode(&updatedWork); err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for i, xx := range y {
		if xx.ID == id {
			updatedWork.ID = id
			y[i] = updatedWork

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updatedWork)
			return
		}
	}

	http.Error(w, "todo not found", http.StatusNotFound)
}

func deletList(w http.ResponseWriter, r *http.Request, id string) {
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for i, xx := range y {
		if xx.ID == id {
			y = append(y[:i], y[i+1:]...)

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
		handlerGet(w, r)
	case "POST":
		createHandler(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func todoByIDHandler(w http.ResponseWriter, r *http.Request) {
	id:=strings.TrimPrefix(r.URL.Path, "/todos/")

	switch r.Method {
	case "PUT":
		putWork(w, r, id)
	case "DELETE":
		deletList(w, r, id)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/todos", todosHandler)
	http.HandleFunc("/todos/", todoByIDHandler)
	http.ListenAndServe(":8080", nil)
}
