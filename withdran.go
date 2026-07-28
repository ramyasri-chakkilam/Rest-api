package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type Account2 struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Balance int    `json:"balance"`
}

var accounts2 = []Account2{
	{ID: 1, Name: "Ramya", Balance: 5000},
	{ID: 2, Name: "Ravi", Balance: 3000},
}

type DepositRequest2 struct {
	Amount int `json:"amount"`
}

func withdraw(w http.ResponseWriter, r *http.Request) {

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

	var req DepositRequest2

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	if req.Amount <= 0 {
		http.Error(w, "Amount should be greater than zero", http.StatusBadRequest)
		return
	}

	for i := range accounts2 {

		if accounts2[i].ID == id {

			if accounts2[i].Balance < req.Amount {
				http.Error(w, "Insufficient Balance", http.StatusBadRequest)
				return
			}

			accounts2[i].Balance -= req.Amount

			err = json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "Amount Withdrawn Successfully",
				"account": accounts2[i],
			})

			if err != nil {
				http.Error(w, "Unable to write response", http.StatusInternalServerError)
				return
			}

			return
		}
	}

	http.Error(w, "Account Not Found", http.StatusNotFound)
}
