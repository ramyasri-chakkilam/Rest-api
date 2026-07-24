package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type Book struct {
	ID              int    `json:"id"`
	Title           string `json:"title"`
	Author          string `json:"author"`
	AvailableCopies int    `json:"available_copies"`
}

var books = []Book{
	{ID: 1, Title: "Go Programming", Author: "John", AvailableCopies: 3},
	{ID: 2, Title: "Microservices", Author: "David", AvailableCopies: 1},
}

func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var newBook Book

	err := json.NewDecoder(r.Body).Decode(&newBook)
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	newBook.ID = len(books) + 1
	books = append(books, newBook)

	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(newBook)
	if err != nil {
		http.Error(w, "Unable to write response", http.StatusInternalServerError)
		return
	}
}

func getBookByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
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

	for _, book := range books {
		if book.ID == id {
			err = json.NewEncoder(w).Encode(book)
			if err != nil {
				http.Error(w, "Unable to write response", http.StatusInternalServerError)
			}
			return
		}
	}

	http.Error(w, "Book Not Found", http.StatusNotFound)
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	err := json.NewEncoder(w).Encode(books)
	if err != nil {
		http.Error(w, "Unable to write response", http.StatusInternalServerError)
	}
}

func issueBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPut {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
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

	for i := range books {

		if books[i].ID == id {

			if books[i].AvailableCopies == 0 {
				http.Error(w, "Book Out of Stock", http.StatusBadRequest)
				return
			}

			books[i].AvailableCopies--

			err = json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "Book issued successfully",
				"book":    books[i],
			})

			if err != nil {
				http.Error(w, "Unable to write response", http.StatusInternalServerError)
				return
			}

			return
		}
	}

	http.Error(w, "Book Not Found", http.StatusNotFound)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodDelete {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
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

	for i := range books {

		if books[i].ID == id {

			books = append(books[:i], books[i+1:]...)

			json.NewEncoder(w).Encode(map[string]string{
				"message": "Book Deleted Successfully",
			})

			return
		}
	}

	http.Error(w, "Book Not Found", http.StatusNotFound)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPut {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
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

	var updatedBook Book

	err = json.NewDecoder(r.Body).Decode(&updatedBook)
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	for i := range books {

		if books[i].ID == id {

			updatedBook.ID = id
			books[i] = updatedBook

			err = json.NewEncoder(w).Encode(books[i])
			if err != nil {
				http.Error(w, "Unable to write response", http.StatusInternalServerError)
			}
			return
		}
	}

	http.Error(w, "Book Not Found", http.StatusNotFound)
}

func returnBook(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPut {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
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

	for i := range books {

		if books[i].ID == id {

			books[i].AvailableCopies++

			err = json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "Book returned successfully",
				"book":    books[i],
			})

			if err != nil {
				http.Error(w, "Unable to write response", http.StatusInternalServerError)
				return
			}

			return
		}
	}

	http.Error(w, "Book Not Found", http.StatusNotFound)
}

func main3() {

	http.HandleFunc("/book", func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {

		case http.MethodPost:
			createBook(w, r)

		case http.MethodGet:
			id := r.URL.Query().Get("id")

			if id == "" {
				getBooks(w, r)
			} else {
				getBookByID(w, r)
			}

		case http.MethodPut:
			updateBook(w, r)

		case http.MethodDelete:
			deleteBook(w, r)

		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/book/issue", issueBook)

	http.HandleFunc("/book/return", returnBook)

	fmt.Println("Server running on :8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
