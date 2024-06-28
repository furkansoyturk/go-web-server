package main

import (
	"log"
	"os"
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

type IDAutoIncrement struct {
	Id int
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func newDB() (*DB, error) {

	db := DB{
		path: "db.json",
	}

	file, err := os.OpenFile(db.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("err while initializing db connection...")
	}
	log.Printf("%v connected", file.Name())

	return &db, nil
}

// CreateChirp creates a new chirp and saves it to disk
// func (db *DB) CreateChirp(body string) (Chirp, error) {
//
// }
//
// // GetChirps returns all chirps in the database
// func (db *DB) GetChirps() ([]Chirp, error) {
//
// }
//
// // ensureDB creates a new database file if it doesn't exist
// func (db *DB) ensureDB() error {
//
// }
//
// // loadDB reads the database file into memory
// func (db *DB) loadDB() (DBStructure, error) {
//
// }
//
// // writeDB writes the database file to disk
// func (db *DB) writeDB(dbStructure DBStructure) error {
//
// }
