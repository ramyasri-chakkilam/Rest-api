package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type Transaction struct {
	Amount int    `json:"amount"`
	RefID  string `json:"refid"`
	Status string `json:"status"`
	TxnID  string `json:"txnid"`
}

var (
	txnStore = make(map[string]Transaction)
	mu       sync.RWMutex
)

func ReqHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Amount int    `json:"amount"`
		RefID  string `json:"refid"`
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}
	id := fmt.Sprintf("TXN%d", time.Now().UnixNano())
	txn := Transaction{
		Amount: req.Amount,
		RefID:  req.RefID,
		TxnID:  id,
		Status: "processing",
	}
	mu.Lock()
	txnStore[id] = txn
	mu.Unlock()

	go asyncProcessing(id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"TransactionId": id,
		"Status":        txn.Status,
	})
}

func asyncProcessing(id string) {
	time.Sleep(3 * time.Second)
	status := "failed"
	if rand.Intn(100) < 80 {
		status = "success"
	}
	mu.Lock()
	txn := txnStore[id]
	txn.Status = status
	txnStore[id] = txn
	mu.Unlock()
}

func getTxnId(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/transaction/"):]
	if r.Method != "GET" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	mu.RLock()
	defer mu.RUnlock()
	txn, ok := txnStore[id]
	if !ok {
		http.Error(w, "Transaction Not Found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(txn)
}

func getStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	mu.RLock()
	defer mu.RUnlock()
	total, processing, success, failed := 0, 0, 0, 0
	for _, txn := range txnStore {
		total++
		switch txn.Status {
		case "processing":
			processing++
		case "success":
			success++
		case "failed":
			failed++

		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{
		"total":      total,
		"processing": processing,
		"success":    success,
		"failed":     failed,
	})
}

func main() {
	http.HandleFunc("/transaction/create", ReqHandler)
	http.HandleFunc("/transaction/stats", getStats)
	http.HandleFunc("/transaction/", getTxnId) // for /transaction/:id
	http.ListenAndServe(":8080", nil)

}
