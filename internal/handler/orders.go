package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/Numbone/practice0/internal/cache"
	"github.com/Numbone/practice0/internal/db"
)

func OrdersHandler(cacheLayer *cache.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orders := cacheLayer.GetAll()

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(orders); err != nil {
			log.Printf("failed to encode orders response: %v\n", err)
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	}
}

func OrderHandler(cacheLayer *cache.Cache, dbConn *db.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orderID := strings.TrimPrefix(r.URL.Path, "/order/")
		if orderID == "" {
			http.Error(w, "order id is required", http.StatusBadRequest)
			return
		}

		if order, ok := cacheLayer.Get(orderID); ok {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(order)
			return
		}

		order, err := dbConn.GetOrder(orderID)
		log.Println(order, err, "result")
		if err != nil {
			http.Error(w, "order not found", http.StatusNotFound)
			return
		}

		cacheLayer.Set(order)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(order)
	}
}
