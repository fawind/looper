package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

// Note is the model for a note
type Note struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
}

// NewNote is the model for a note to create
type NewNote struct {
	Content string `json:"content"`
}

// NoteID is the model for a new ID response
type NoteID struct {
	ID int `json:"id"`
}

var notesIDServiceEndpoint string
var notesStoreServiceEndpoint string

func main() {
	port := getEnv("PORT", "8080")
	host := getEnv("HOST", "0.0.0.0")
	notesIDServiceEndpoint = "http://" + getRequiredEnv("NOTES_ID_SERVICE")
	notesStoreServiceEndpoint = "http://" + getRequiredEnv("NOTES_STORE_SERVICE")

	router := mux.NewRouter()

	router.HandleFunc("/", HandleIndex)
	router.HandleFunc("/notes", HandleGetNotes).Methods("GET")
	router.HandleFunc("/notes", HandleAddNote).Methods("POST")

	log.Printf("Server running on port %s:%s\n", host, port)
	log.Printf("Notes-ID Service Endpoint: %s\n", notesIDServiceEndpoint)
	log.Printf("Notes-Store Service Endpoint: %s\n", notesStoreServiceEndpoint)
	log.Fatal(http.ListenAndServe(host+":"+port, router))
}

// HandleIndex returns the service name
func HandleIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Notes Service")
}

// HandleGetNotes returns all stored notes
func HandleGetNotes(w http.ResponseWriter, r *http.Request) {
	notes, err := fetchNotes()
	if err != nil {
		log.Printf("Error fetching notes: %s\n", err)
		http.Error(w, "Error fetching notes", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(notes)
}

// HandleAddNote adds a new note to the store
func HandleAddNote(w http.ResponseWriter, r *http.Request) {
	var newNote NewNote
	err := json.NewDecoder(r.Body).Decode(&newNote)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Printf("Error decoding request body: %s\n", err)
		return
	}
	noteID, err := getNextID()
	if err != nil {
		http.Error(w, "Error creating the new note", http.StatusBadRequest)
		log.Printf("Error fetching next note id: %s\n", err)
		return
	}
	note := Note{ID: noteID, Content: newNote.Content}
	err = addNote(note)
	if err != nil {
		http.Error(w, "Error creating the new note", http.StatusBadRequest)
		log.Printf("Error adding new note to store: %s\n", err)
	}
	json.NewEncoder(w).Encode(note)
}

func fetchNotes() ([]Note, error) {
	r, err := http.Get(notesStoreServiceEndpoint + "/notes")
	if err != nil {
		return nil, err
	}
	var notes []Note
	err = json.NewDecoder(r.Body).Decode(&notes)
	if err != nil {
		return nil, err
	}
	return notes, nil
}

func addNote(note Note) error {
	body, err := json.Marshal(note)
	if err != nil {
		return err
	}
	_, err = http.Post(
		notesStoreServiceEndpoint+"/notes",
		"application/json",
		bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	return nil
}

func getNextID() (int, error) {
	r, err := http.Get(notesIDServiceEndpoint + "/next-id")
	if err != nil {
		return -1, err
	}
	var nextID NoteID
	err = json.NewDecoder(r.Body).Decode(&nextID)
	if err != nil {
		return -1, err
	}
	return nextID.ID, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getRequiredEnv(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	log.Fatalf("Env variable %s not set\n", key)
	return ""
}
