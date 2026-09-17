package main

import (
	"encoding/json"
	"net/http"
)

func main() {
	client, err := NewClient()
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "POST required", http.StatusMethodNotAllowed)
			return
		}
		var in struct {
			Query string `json:"query"`
		}
		if json.NewDecoder(r.Body).Decode(&in) != nil || in.Query == "" {
			http.Error(w, "query is required", http.StatusBadRequest)
			return
		}
		result, err := client.Search(r.Context(), in.Query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
