package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type Stock struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
}

var stocks = []Stock{
	{ID: 1, Name: "HP", Quantity: 3},
	{ID: 2, Name: "Dell", Quantity: 5},
}

func createStock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var newStock Stock

	err := json.NewDecoder(r.Body).Decode(&newStock)
	if err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	newStock.ID = len(stocks) + 1
	stocks = append(stocks, newStock)

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(map[string]string{
		"message": "Successfully Created",
	})
	if err != nil {
		http.Error(w, "Unable to write response", http.StatusInternalServerError)
		return
	}
}

func deleteStockByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Allow only DELETE method
	if r.Method != http.MethodDelete {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get ID
	idStr := strings.TrimPrefix(r.URL.Path, "/stock/")
	if idStr == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	// Get count (default = 1)
	count := 1
	countStr := r.URL.Query().Get("count")

	if countStr != "" {
		count, err = strconv.Atoi(countStr)
		if err != nil {
			http.Error(w, "Invalid count", http.StatusBadRequest)
			return
		}
	}

	// Count must be greater than 0
	if count <= 0 {
		http.Error(w, "Count must be greater than 0", http.StatusBadRequest)
		return
	}

	// Find stock
	for i := range stocks {

		if stocks[i].ID == id {

			// No stock available
			if stocks[i].Quantity == 0 {
				http.Error(w, "Stock is already empty", http.StatusBadRequest)
				return
			}

			// Cannot remove more than available
			if count > stocks[i].Quantity {
				http.Error(w, "Requested quantity exceeds available stock", http.StatusBadRequest)
				return
			}

			// Reduce quantity
			stocks[i].Quantity -= count

			// Optional:
			// Remove stock record if quantity becomes 0
			/*
				if stocks[i].Quantity == 0 {
					stocks = append(stocks[:i], stocks[i+1:]...)
				}
			*/

			err = json.NewEncoder(w).Encode(map[string]interface{}{
				"message":           "Stock updated successfully",
				"removedQuantity":   count,
				"remainingQuantity": stocks[i].Quantity,
			})
			if err != nil {
				http.Error(w, "Unable to write response", http.StatusInternalServerError)
			}
			return
		}
	}

	// ID not found
	http.Error(w, "Stock not found", http.StatusNotFound)
}

func getStockByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/stock/")
	if idStr == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}        
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if id == 0 {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	for _, item := range stocks {
		if item.ID == id {
			err = json.NewEncoder(w).Encode(item)
			if err != nil {
				http.Error(w, "Unable to write response", http.StatusInternalServerError)
			}
			return
		}
	}

	http.Error(w, "Stock Not Found", http.StatusNotFound)
}

func getAllStocks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	err := json.NewEncoder(w).Encode(stocks)
	if err != nil {
		http.Error(w, "Unable to write response", http.StatusInternalServerError)
		return
	}
}

func main2() {
	http.HandleFunc("/stock", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/stock" {
			http.NotFound(w, r)
			return
		}

		switch r.Method {
		case http.MethodPost:
			createStock(w, r)
		case http.MethodGet:
			getAllStocks(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/stock/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodDelete:
			deleteStockByID(w, r)
		case http.MethodGet:
			getStockByID(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	http.ListenAndServe(":8080", nil)
}


