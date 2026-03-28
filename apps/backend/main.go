package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type HelloResponse struct {
	Message string `json:"message"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Optional: serve the Expo web build from a static directory.
	// Set STATIC_DIR to the path of apps/client/dist to serve the frontend
	// from the same origin as the API.
	staticDir := os.Getenv("STATIC_DIR")

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		resp := HelloResponse{
			Message: "Hello from the Go backend! Foreningsliv is alive 🎉",
		}
		json.NewEncoder(w).Encode(resp)
	})

	// Handle CORS preflight
	mux.HandleFunc("OPTIONS /api/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.WriteHeader(http.StatusNoContent)
	})

	// Serve static files if STATIC_DIR is set
	if staticDir != "" {
		fs := http.FileServer(http.Dir(staticDir))
		mux.Handle("/", fs)
		fmt.Printf("Serving static files from %s\n", staticDir)
	}

	fmt.Printf("Backend listening on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
