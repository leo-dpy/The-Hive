package api

import (
	"encoding/json"
	"net/http"
)

// PingHandler is a simple endpoint to test backend connectivity.
func PingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	
	json.NewEncoder(w).Encode(map[string]string{
		"message": "pong from The Hive API",
		"status":  "ok",
		"engine":  "Golang",
	})
}
