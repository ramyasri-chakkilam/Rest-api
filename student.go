package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type Student struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Marks int    `json:"marks"`
}

var students = []Student{
	{Name: "Ramya", Age: 26, Marks: 98},
	{Name: "Teja", Age: 24, Marks: 68},
}

func putStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPut {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	var newStudent Student
	err := json.NewDecoder(r.Body).Decode(&newStudent)
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	for i := range students {
		if students[i].ID == id {
			newStudent.ID = id
			students[i] = newStudent

			w.WriteHeader(http.StatusOK)
			err = json.NewEncoder(w).Encode(newStudent)
			if err != nil {
				http.Error(w, "Unable to write response", http.StatusInternalServerError)
				return
			}
			return
		}
	}
	http.Error(w, "student not found", http.StatusNotFound)
}
