package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type Employee struct {
	Name   string `json:"name"`
	ID     int    `json:"id"`
	Age    int    `json:"age"`
	Salary int    `json:"salary"`
}

var emp = []Employee{
	{ID: 1, Name: "Ram", Salary: 5000},
	{ID: 2, Name: "Sita", Salary: 6000},
}

func createEmp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}
	var newEmp Employee
	err := json.NewDecoder(r.Body).Decode(&newEmp)
	if err != nil {
		http.Error(w, "Invalid Request payload", http.StatusBadRequest)
		return
	}
	newEmp.ID = len(emp) + 1
	emp = append(emp, newEmp)
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(map[string]string{
		"message": "successfully created",
	})
	if err != nil {
		http.Error(w, "unable to write response", http.StatusInternalServerError)
		return
	}
}

func getEmp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	err := json.NewEncoder(w).Encode(emp)
	if err != nil {
		http.Error(w, "Unable to write response", http.StatusInternalServerError)
		return
	}
}
func getEmpByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
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
	for _, emp1 := range emp {
		if emp1.ID == id {
			err = json.NewEncoder(w).Encode(emp1)
			if err != nil {
				http.Error(w, "Unable to write response", http.StatusInternalServerError)
			}
			return
		}
	}

	http.Error(w, "Stock Not Found", http.StatusNotFound)
}
func delEmp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
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
	for i, emps := range emp {
		if emps.ID == id {
			emp = append(emp[:i], emp[i+1:]...)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "emp Deleted Successfully",
			})
			return
		}
	}
	http.Error(w, "Stock Not Found", http.StatusNotFound)
}

func mai1n() {
	http.HandleFunc("/emp", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {

		case http.MethodPost:
			createEmp(w, r)

		case http.MethodGet:
			id := r.URL.Query().Get("id")
			if id == "" {
				getEmp(w, r)
			} else {
				getEmpByID(w, r)
			}

		case http.MethodDelete:

			delEmp(w, r)

		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}

	})
	http.ListenAndServe(":8080", nil)
}
