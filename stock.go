package main

// import (
// 	"encoding/json"
// 	"net/http"
// 	"strconv"
// 	"strings"
// )

// type Stock struct {
// 	ID       int    `json:"id"`
// 	Name     string `json:"name"`
// 	Quantity int    `json:"quantity"`
// }

// var stocks = []Stock{
// 	{1, "Dell", 10},
// 	{2, "HP", 5},
// 	{3, "Lenovo", 8},
// }

// // --------------------
// // POST /stock
// // --------------------

// func createStock(w http.ResponseWriter, r *http.Request) {

// 	if r.Method != http.MethodPost {
// 		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	var stock Stock

// 	err := json.NewDecoder(r.Body).Decode(&stock)
// 	if err != nil {
// 		http.Error(w, "Invalid Request", http.StatusBadRequest)
// 		return
// 	}

// 	stock.ID = len(stocks) + 1

// 	stocks = append(stocks, stock)

// 	w.Header().Set("Content-Type", "application/json")

// 	json.NewEncoder(w).Encode(map[string]string{
// 		"message": "Stock Created Successfully",
// 	})
// }

// // --------------------
// // GET /stock
// // --------------------

// func getAllStocks(w http.ResponseWriter, r *http.Request) {

// 	if r.Method != http.MethodGet {
// 		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")

// 	json.NewEncoder(w).Encode(stocks)
// }

// // --------------------
// // GET /stock/{id}
// // --------------------

// func getStockByID(w http.ResponseWriter, r *http.Request) {

// 	if r.Method != http.MethodGet {
// 		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	idStr := strings.TrimPrefix(r.URL.Path, "/stock/")

// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		http.Error(w, "Invalid ID", http.StatusBadRequest)
// 		return
// 	}

// 	for _, stock := range stocks {

// 		if stock.ID == id {

// 			w.Header().Set("Content-Type", "application/json")

// 			json.NewEncoder(w).Encode(stock)
// 			return
// 		}
// 	}

// 	http.Error(w, "Stock Not Found", http.StatusNotFound)
// }

// // --------------------
// // DELETE /stock/{id}
// // --------------------

// func deleteStock(w http.ResponseWriter, r *http.Request) {

// 	if r.Method != http.MethodDelete {
// 		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	idStr := strings.TrimPrefix(r.URL.Path, "/stock/")

// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		http.Error(w, "Invalid ID", http.StatusBadRequest)
// 		return
// 	}

// 	for i, stock := range stocks {

// 		if stock.ID == id {

// 			stocks = append(stocks[:i], stocks[i+1:]...)

// 			w.Header().Set("Content-Type", "application/json")

// 			json.NewEncoder(w).Encode(map[string]string{
// 				"message": "Stock Deleted Successfully",
// 			})

// 			return
// 		}
// 	}	

// 	http.Error(w, "Stock Not Found", http.StatusNotFound)
// }

// func main() {

// 	http.HandleFunc("/stock", func(w http.ResponseWriter, r *http.Request) {

// 		if r.URL.Path != "/stock" {
// 			http.NotFound(w, r)
// 			return
// 		}

// 		switch r.Method {

// 		case http.MethodPost:
// 			createStock(w, r)

// 		case http.MethodGet:
// 			getAllStocks(w, r)

// 		default:
// 			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
// 		}
// 	})

// 	http.HandleFunc("/stock/", func(w http.ResponseWriter, r *http.Request) {

// 		switch r.Method {

// 		case http.MethodGet:
// 			getStockByID(w, r)

// 		case http.MethodDelete:
// 			deleteStock(w, r)

// 		default:
// 			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
// 		}
// 	})

// 	http.ListenAndServe(":8080", nil)
// }