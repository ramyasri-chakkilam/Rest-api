package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type Order struct {
	ID     int    `json:"id"`
	Item   string `json:"item"`
	Status string `json:"status"`
}

type OrderRequest struct {
	Item string `json:"item"`
}

type StatusReq struct {
	Status string `json:"status"`
}

// var orders []Order
var orders = []Order{
	{
		ID:     1,
		Item:   "Laptop",
		Status: "Pending",
	},
	{
		ID:     2,
		Item:   "Mobile",
		Status: "Processing",
	},
}

func updateStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPut {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	var req StatusReq
	idStr := strings.TrimPrefix(r.URL.Path, "/order/")
	if idStr == "" {
		http.Error(w, "ID Is Required", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid Request Payload", http.StatusBadRequest)
		return
	}
	if req.Status == "" {
		http.Error(w, "Status Is Required", http.StatusBadRequest)
		return
	}
	for i := range orders {
		if orders[i].ID == id {
			switch orders[i].Status {
			case "Pending":
				if req.Status != "Processing" {
					http.Error(w, "Cannot Skip Order Status", http.StatusBadRequest)
					return
				}
			case "Processing":
				if req.Status != "Shipped" {
					http.Error(w, "Cannot Skip Order Status", http.StatusBadRequest)
					return
				}
			case "Shipped":
				if req.Status != "Delivered" {
					http.Error(w, "Cannot Skip Order Status", http.StatusBadRequest)
					return
				}
			case "Delivered":

				http.Error(w, "Item Already Delivered", http.StatusBadRequest)
				return
			}
			orders[i].Status = req.Status
			err = json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "succesfully updated the order",
				"order":   orders[i],
			})
			if err != nil {
				http.Error(w, "Error Writing the Response", http.StatusBadRequest)
				return
			}
			return
		}
	}
	http.Error(w, "Order Not Found", http.StatusNotFound)

}

func getOrders(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	err := json.NewEncoder(w).Encode(orders)
	if err != nil {
		http.Error(w, "Error Writing Response", http.StatusInternalServerError)
		return
	}
}

func createOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req OrderRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid Request Payload", http.StatusBadRequest)
		return
	}

	if req.Item == "" {
		http.Error(w, "Item Is Required", http.StatusBadRequest)
		return
	}

	order := Order{
		ID:     len(orders) + 1,
		Item:   req.Item,
		Status: "Pending",
	}

	orders = append(orders, order)

	err = json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Order Created Successfully",
		"order":   order,
	})
	if err != nil {
		http.Error(w, "Error Writing Response", http.StatusInternalServerError)
		return
	}
}

func delOrder(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodDelete {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/order/")

	if idStr == "" {
		http.Error(w, "ID Is Required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	for i := range orders {

		if orders[i].ID == id {

			orders = append(orders[:i], orders[i+1:]...)

			err = json.NewEncoder(w).Encode(map[string]interface{}{
				"message": "Order Deleted Successfully",
			})
			if err != nil {
				http.Error(w, "Error Writing Response", http.StatusInternalServerError)
				return
			}

			return
		}
	}

	http.Error(w, "Order Not Found", http.StatusNotFound)
}

func main() {

	http.HandleFunc("/order/", func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {

		case http.MethodPut:
			updateStatus(w, r)

		case http.MethodDelete:
			delOrder(w, r)

		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/order", func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {

		case http.MethodPost:
			createOrder(w, r)

		case http.MethodGet:
			getOrders(w, r)

		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	http.ListenAndServe(":8080", nil)
}
