package main

import (
	"encoding/json"
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

func writeJSONError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, ErrorResponse{
		Status: status,
		Text:   http.StatusText(status),
		Error:  message,
	})
}

// Mensjae de error en la validacion de datos
func writeJSONValidationErrors(w http.ResponseWriter, errors []string) {
	writeJSON(w, http.StatusBadRequest, ErrorResponse{
		Status: http.StatusBadRequest,
		Text:   http.StatusText(http.StatusBadRequest),
		Errors: errors,
	})
}

// Validacion de datos de entrada para modificar o crear una serie
func validateSeriesInput(s Series) []string {
	var errors []string

	if strings.TrimSpace(s.Name) == "" {
		errors = append(errors, "name is required")
	}

	if s.TotalEpisodes <= 0 {
		errors = append(errors, "total_episodes must be greater than 0")
	}

	if s.CurrentEpisode < 0 {
		errors = append(errors, "current_episode cannot be negative")
	}

	if s.TotalEpisodes > 0 && s.CurrentEpisode > s.TotalEpisodes {
		errors = append(errors, "current_episode cannot be greater than total_episodes")
	}

	return errors
}

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

// Funcion para extraer el ID de la ruta para poder buscar la entrada
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
