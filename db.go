package main

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func initDB() {
	var err error
	
	// Obtener ruta absoluta de la base de datos
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal("Error getting working directory:", err)
	}
	dbPath := filepath.Join(wd, "series2.db")
	log.Println("Database path:", dbPath)
	
	db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
}
