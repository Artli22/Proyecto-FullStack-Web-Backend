package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

var db *sql.DB

func initDB() {
	dbURL := os.Getenv("TURSO_DATABASE_URL") 
	authToken := os.Getenv("TURSO_AUTH_TOKEN")

	if dbURL == "" {
		log.Fatal("TURSO_DATABASE_URL no está definida")
	}
	if authToken == "" {
		log.Fatal("TURSO_AUTH_TOKEN no está definida")
	}

	// El driver espera la URL con el token como query param
	connStr := fmt.Sprintf("%s?authToken=%s", dbURL, authToken)

	var err error
	db, err = sql.Open("libsql", connStr)
	if err != nil {
		log.Fatal("Error abriendo conexión a Turso:", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("Error conectando a Turso:", err)
	}

	log.Println("Conectado a Turso:", dbURL)
}