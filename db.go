package main

import (
	"sync"
)

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {

	err := os.WriteFile("testdata/hello", []byte("Hello, Gophers!"), 0666)
	if err != nil {
		log.Fatal(err)
	}
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error) {

}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {

}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {

}

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error) {

}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {

}
