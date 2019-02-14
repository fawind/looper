package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"sync"
)

var notes = []Note{}
var noteLock sync.Mutex

// Note is the model for a note
type Note struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
}

func main() {
	port := getEnv("PORT", "8080")
	host := getEnv("HOST", "0.0.0.0")

	router := mux.NewRouter()

	router.HandleFunc("/", HandleIndex)
	router.HandleFunc("/notes", HandleGetNotes).Methods("GET")
	router.HandleFunc("/notes", HandleAddNote).Methods("POST")

	log.Printf("Server running on port %s:%s\n", host, port)
	log.Fatal(http.ListenAndServe(host+":"+port, router))
}

// HandleIndex returns the service name
func HandleIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Notes Store Service")
}

// HandleGetNotes returns all stored notes
func HandleGetNotes(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(notes)
}

// HandleAddNote adds a new note to the store
func HandleAddNote(w http.ResponseWriter, r *http.Request) {
	var note Note
	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Printf("Error decoding request body: %s\n", err)
		return
	}
	noteLock.Lock()
	defer noteLock.Unlock()
	if isDuplicateID(note.ID) {
		http.Error(w, "Note ID already exists", http.StatusBadRequest)
	} else {
		notes = append(notes, note)
	}
}

func isDuplicateID(id int) bool {
	for _, note := range notes {
		if note.ID == id {
			return true
		}
	}
	return false
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
