package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
)

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

type DBConnection struct {
	path  string
	mux   *sync.RWMutex
	index int
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB() (*DBConnection, error) {

	db := DBConnection{
		path: "db.json",
	}

	file, err := os.OpenFile(db.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("err while initializing db connection...")
	}
	defer file.Close()
	log.Printf("%v connected", file.Name())
	data, err := os.ReadFile(file.Name())
	if err != nil {
		log.Panicf("failed reading data from file: %s", err)
	}

	var chirpMap DBStructure
	json.Unmarshal(data, &chirpMap)
	log.Printf("len of chirp list %v", len(chirpMap.Chirps))
	// db.getDBStructure()
	return &db, nil
}

// CreateChirp creates a new chirp and saves it to disk
func (db *DBConnection) save(body string) (Chirp, error) {
	db.index = len(db.getDBStructure().Chirps) + 1
	fmt.Printf("current index of db -> %v", db.index)
	fmt.Println("----------")

	chirp := Chirp{
		Id:   db.index,
		Body: body,
	}
	var m = make(map[int]Chirp)

	fmt.Println("empty map initialized")
	if db.index > 1 {
		m = db.getDBStructure().Chirps
		m[db.index] = chirp

	} else {
		fmt.Println("there is no key in map log-")
		m[1] = chirp
		fmt.Println("error while adding chrip to map")
	}
	dbStructure := DBStructure{
		Chirps: m,
	}
	data, _ := json.Marshal(dbStructure)
	os.WriteFile(db.path, []byte(data), 0644)

	log.Println("writed successfully")
	return chirp, nil
}

func (db *DBConnection) getDBStructure() DBStructure {

	file, err := os.OpenFile(db.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("err while initializing db connection...")
	}
	defer file.Close()
	log.Printf("%v connected", file.Name())
	data, err := os.ReadFile(file.Name())
	if err != nil {
		log.Panicf("failed reading data from file: %s", err)
	}

	var dbImg DBStructure
	json.Unmarshal(data, &dbImg)
	log.Println("DB")
	log.Println(dbImg)
	log.Printf("len of chirp list %v", len(dbImg.Chirps))
	log.Printf("last id of created chirp %v", len(dbImg.Chirps))

	return dbImg
}

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
