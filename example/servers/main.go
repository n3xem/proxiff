package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"time"
)

type Response struct {
	Message   string    `json:"message"`
	Version   string    `json:"version"`
	Timestamp time.Time `json:"timestamp"`
}

func main() {
	port := flag.String("port", "8081", "Port to listen on")
	version := flag.String("version", "v1.0", "Server version")
	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handleRequest(w, r, *version)
	})

	http.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		handleUsers(w, r, *version)
	})

	http.HandleFunc("/api/status", func(w http.ResponseWriter, r *http.Request) {
		handleStatus(w, r, *version)
	})

	addr := ":" + *port
	log.Printf("Starting sample server %s on %s", *version, addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request, version string) {
	resp := Response{
		Message:   "Hello from " + version,
		Version:   version,
		Timestamp: time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func handleUsers(w http.ResponseWriter, r *http.Request, version string) {
	users := []map[string]interface{}{
		{"id": 1, "name": "Alice"},
		{"id": 2, "name": "Bob"},
	}

	// Newer version has an additional field
	if version == "newer" || version == "v2.0" {
		users = []map[string]interface{}{
			{"id": 1, "name": "Alice", "email": "alice@example.com"},
			{"id": 2, "name": "Bob", "email": "bob@example.com"},
		}
	}

	resp := map[string]interface{}{
		"users":   users,
		"version": version,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func handleStatus(w http.ResponseWriter, r *http.Request, version string) {
	status := map[string]interface{}{
		"status":  "ok",
		"version": version,
	}

	// Current version has different status code for demonstration
	if version == "current" || version == "v1.0" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	} else {
		// Newer version returns 201 instead of 200
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
	}

	json.NewEncoder(w).Encode(status)
}
