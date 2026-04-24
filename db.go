package main

import (
    "database/sql"
    "errors"
    "fmt"
    "os"
    "sync"

    _ "github.com/tursodatabase/libsql-client-go/libsql"
)

var db *sql.DB

var (
    dbInitOnce sync.Once
    dbInitErr  error
)

func initDB() error {
    dbURL := os.Getenv("TURSO_DATABASE_URL")
    authToken := os.Getenv("TURSO_AUTH_TOKEN")

    if dbURL == "" {
        return errors.New("TURSO_DATABASE_URL no esta definida")
    }
    if authToken == "" {
        return errors.New("TURSO_AUTH_TOKEN no esta definida")
    }

    // El driver espera la URL con el token como query param
    connStr := fmt.Sprintf("%s?authToken=%s", dbURL, authToken)

    var err error
    db, err = sql.Open("libsql", connStr)
    if err != nil {
        return fmt.Errorf("error abriendo conexion a Turso: %w", err)
    }

    if err = db.Ping(); err != nil {
        return fmt.Errorf("error conectando a Turso: %w", err)
    }

    return nil
}

func ensureDB() error {
    dbInitOnce.Do(func() {
        dbInitErr = initDB()
    })

    return dbInitErr
}