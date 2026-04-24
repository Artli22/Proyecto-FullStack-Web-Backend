package main
// Archivo helpers.go

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type ErrorResponse struct {
	Status int      `json:"status"`
	Text   string   `json:"text"`
	Error  string   `json:"error,omitempty"`
	Errors []string `json:"errors,omitempty"`
}

// Funcion para escribir respuestas JSON y errores de validacion
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, `{"status":500,"text":"Internal Server Error","error":"could not encode response"}`, http.StatusInternalServerError)
	}
}

// Mensaje de error simple
func writeJSONError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, ErrorResponse{
		Status: status,
		Text:   http.StatusText(status),
		Error:  message,
	})
}

// Mensjae de lista de errores; si se detectan multiples
func writeJSONValidationErrors(w http.ResponseWriter, errors []string) {
	writeJSON(w, http.StatusBadRequest, ErrorResponse{
		Status: http.StatusBadRequest,
		Text:   http.StatusText(http.StatusBadRequest),
		Errors: errors,
	})
}

func ensureDBOrWriteError(w http.ResponseWriter) bool {
	if err := ensureDB(); err != nil {
		log.Printf("database initialization failed: %v", err)
		writeJSONError(w, http.StatusInternalServerError, "database is not available")
		return false
	}

	return true
}

// Validacion de datos de entrada para modificar o crear una serie
func validateSeriesInput(s Series) []string {
	var errors []string

	if strings.TrimSpace(s.Name) == "" {
		errors = append(errors, "el nombre es obligatorio")
	}

	if s.TotalEpisodes <= 0 {
		errors = append(errors, "el numero de total_episodes debe ser mayor a 0")
	}

	if s.CurrentEpisode < 0 {
		errors = append(errors, "current_episode no puede ser negativo")
	}

	if s.TotalEpisodes > 0 && s.CurrentEpisode > s.TotalEpisodes {
		errors = append(errors, "current_episode no puede ser mayor a total_episodes")
	}

	return errors
}

// Conversion de texto a un entero positivo
func parsePositiveInt(value string, defaultValue int) int {
	n, err := strconv.Atoi(value)
	if err != nil || n < 1 {
		return defaultValue
	}
	return n
}

// Funcion para obtener los parametros necesarios para la paginacion de las entradas
func getPaginationParams(r *http.Request) (int, int, int) {
	page := parsePositiveInt(r.URL.Query().Get("page"), 1)
	limit := parsePositiveInt(r.URL.Query().Get("limit"), 10)
	offset := (page - 1) * limit
	return page, limit, offset
}

// Funcion para extraer el parametro de busqueda de la ruta 
func getSearchParam(r *http.Request) string {
	return strings.TrimSpace(r.URL.Query().Get("q"))
}

// Funcion para ordenar las entradas segun el orden ascendente o descendente
func getSortParams(r *http.Request) (string, string) {
	sort := r.URL.Query().Get("sort")
	order := strings.ToLower(r.URL.Query().Get("order"))

	allowedSorts := map[string]bool{
		"id":              true,
		"name":            true,
		"current_episode": true,
		"total_episodes":  true,
	}

	if !allowedSorts[sort] {
		sort = "id"
	}

	if order != "desc" {
		order = "asc"
	}

	return sort, order
}