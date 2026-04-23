package main
// Archivo Main.go

import (
	"log"
	"net/http"
)

// Habilitar CORS para permitir la solicitudes desde el frontend
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func seriesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getSeriesHandler(w, r)
	case http.MethodPost:
		createSeriesHandler(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// Handler para rutas con ID (GET, PUT, DELETE)
func seriesByIDHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getSeriesByIDHandler(w, r)
	case http.MethodPut:
		updateSeriesHandler(w, r)
	case http.MethodDelete:
		deleteSeriesHandler(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// Funcion para iniciar el backend en el puerto 8080
func main() {
	initDB()

	mux := http.NewServeMux()
	mux.HandleFunc("/series", seriesHandler)
	mux.HandleFunc("/series/", seriesByIDHandler)

	log.Println("Listening on :8080")
	if err := http.ListenAndServe(":8080", enableCORS(mux)); err != nil {
		log.Fatal(err)
	}
}
