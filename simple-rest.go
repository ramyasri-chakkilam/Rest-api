package main

import (
	"encoding/json"
	"net/http"
)

type User struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email"`
}

var users = []User{
	{Name: "ramya", Age: 30, Email: "ramya@example.com"},
	{Name: "teja", Age: 25, Email: "teja@example.com"},
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(users)
	return
	}
	if r.Method == "POST" {
	var item User
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	
	if item.Name == "" || item.Email == "" || item.Age <= 0 {
		http.Error(w, "Missing Required Fields", http.StatusBadRequest)
		return
	}
	users = append(users, item)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
	return
	}
	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}

// func main() {
// 	http.HandleFunc("/users", handler)
// 	http.ListenAndServe(":8080", nil)
// }
