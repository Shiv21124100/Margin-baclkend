package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
)

type Asset struct {
	Symbol          string    `json:"symbol"`
	MarkPrice       float64   `json:"mark_price"`
	ContractValue   float64   `json:"contract_value"`
	AllowedLeverage []float64 `json:"allowed_leverage"`
}

type MarginRequest struct {
	Asset        string  `json:"asset"`
	OrderSize    float64 `json:"order_size"`
	Side         string  `json:"side"`
	Leverage     float64 `json:"leverage"`
	MarginClient float64 `json:"margin_client"`
}

type MarginResponse struct {
	Status         string  `json:"status"`
	MarginRequired float64 `json:"margin_required"`
}

type ErrorResponse struct {
	Status         string  `json:"status"`
	Message        string  `json:"message"`
	MarginRequired float64 `json:"margin_required"`
}

var assets = []Asset{
	{
		Symbol:          "BTC",
		MarkPrice:       62000,
		ContractValue:   0.001,
		AllowedLeverage: []float64{5, 10, 20, 50, 100},
	},
	{
		Symbol:          "ETH",
		MarkPrice:       3200,
		ContractValue:   0.01,
		AllowedLeverage: []float64{5, 10, 25, 50},
	},
}

func main() {
	http.HandleFunc("/config/assets", corsMiddleware(configHandler))
	http.HandleFunc("/margin/validate", corsMiddleware(validateHandler))

	fmt.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// corsMiddleware handles CORS headers
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			return
		}

		next(w, r)
	}
}

func configHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string][]Asset{"assets": assets})
}

func validateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req MarginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Find asset
	var selectedAsset *Asset
	for _, a := range assets {
		if a.Symbol == req.Asset {
			selectedAsset = &a
			break
		}
	}

	if selectedAsset == nil {
		http.Error(w, "Asset not found", http.StatusBadRequest)
		return
	}

	// Validate leverage
	isValidLeverage := false
	for _, l := range selectedAsset.AllowedLeverage {
		if l == req.Leverage {
			isValidLeverage = true
			break
		}
	}
	if !isValidLeverage {
		http.Error(w, "Invalid leverage", http.StatusBadRequest)
		return
	}

	// Calculate required margin
	// Formula: (mark_price * order_size * contract_value) / leverage
	marginRequired := (selectedAsset.MarkPrice * req.OrderSize * selectedAsset.ContractValue) / req.Leverage
	
	// Round to 2 decimal places for comparison
	marginRequired = math.Round(marginRequired*100) / 100

	w.Header().Set("Content-Type", "application/json")

	// Compare with client margin (allowing for small float differences if needed, but rounding should handle it)
	if req.MarginClient < marginRequired {
		resp := ErrorResponse{
			Status:         "error",
			Message:        "Insufficient margin submitted",
			MarginRequired: marginRequired,
		}
		w.WriteHeader(http.StatusBadRequest) // Or 200 with error status, but assignment implies rejection
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp := MarginResponse{
		Status:         "ok",
		MarginRequired: marginRequired,
	}
	json.NewEncoder(w).Encode(resp)
}
