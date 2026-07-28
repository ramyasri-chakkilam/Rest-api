package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type Stud struct {
	Name  string `json:"name"`
	Marks int    `json:"marks"`
	ID    int    `json:"id"`
}

type studentReq struct {
	Name  string `json:"name"`
	Marks int    `json:"marks"`
}

type studentResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Marks int    `json:"marks"`
	Grade string `json:"grade"`
}

var studs = []Stud{
	{ID: 1, Name: "Ramya", Marks: 92},
	{ID: 2, Name: "John", Marks: 76},
}

func createStudent(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req studentReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	if req.Marks < 0 || req.Marks > 100 {
		http.Error(w, "Marks should be between 0 and 100", http.StatusBadRequest)
		return
	}

	student := Stud{
		ID:    len(studs) + 1,
		Name:  req.Name,
		Marks: req.Marks,
	}

	studs = append(studs, student)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(student)
}

func getStu(w http.ResponseWriter, r *http.Request) {

	var req studentResponse
	idStr := strings.TrimPrefix(r.URL.Path, "/student/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID is invalid", http.StatusBadRequest)
		return
	}
	//main logic
	for i := range studs {
		if studs[i].ID == id {
			req = studentResponse{
				ID:    studs[i].ID,
				Name:  studs[i].Name,
				Marks: studs[i].Marks,
				Grade: calcalateGrade(studs[i].Marks),
			}
			err = json.NewEncoder(w).Encode(req)
			return
		}
	}
	http.Error(w, "Student not found", http.StatusNotFound)
}

func getAllStudents(w http.ResponseWriter, r *http.Request) {

	// Check HTTP method
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Response slice
	var responses []studentResponse

	// Loop through all students
	for _, student := range studs {

		response := studentResponse{
			ID:    student.ID,
			Name:  student.Name,
			Marks: student.Marks,
			Grade: calcalateGrade(student.Marks),
		}

		responses = append(responses, response)
	}

	// Send response
	if err := json.NewEncoder(w).Encode(responses); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func calcalateGrade(marks int) string {
	if marks >= 90 {
		return "A"
	} else if marks >= 80 {
		return "B"
	} else if marks >= 70 {
		return "C"
	}
	return "D"

}

func main4() {

	http.HandleFunc("/student", func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {

		case http.MethodPost:
			createStudent(w, r)

		case http.MethodGet:
			getAllStudents(w, r)

		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}

	})

	http.HandleFunc("/student/", getStu)

	http.ListenAndServe(":8080", nil)
}
