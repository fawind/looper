package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)

var idCounter int
var idLock sync.Mutex

// IDResponse is the response model for a new ID
type IDResponse struct {
	ID int `json:"id"`
}

func main() {
	port := getEnv("PORT", "8080")
	host := getEnv("HOST", "0.0.0.0")

	http.HandleFunc("/", HandleIndex)
	http.HandleFunc("/next-id", HandleNextID)

	log.Printf("Server running on port %s:%s\n", host, port)
	log.Fatal(http.ListenAndServe(host+":"+port, nil))
}

// HandleIndex returns the service name
func HandleIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Note-ID Service")
}

// HandleNextID returns the next available id
func HandleNextID(w http.ResponseWriter, r *http.Request) {
	nextID := getNextID()
	json.NewEncoder(w).Encode(IDResponse{ID: nextID})
}

func getNextID() int {
	idLock.Lock()
	nextID := idCounter
	idCounter++
	idLock.Unlock()
	return nextID
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
