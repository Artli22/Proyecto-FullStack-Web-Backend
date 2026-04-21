package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// Funcion para extraer el ID desde rutas
func getIDFromPath(path string) (int, error) {
	idStr := strings.TrimPrefix(path, "/series/")
	return strconv.Atoi(idStr)
}

// Handler GET /series
func getSeriesHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
		SELECT id, name, description, image_url, current_episode, total_episodes
		FROM series2
	`)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "could not fetch series")
		return
	}
	defer rows.Close()

	var seriesList []Series

	for rows.Next() {
		var s Series
		err := rows.Scan(
			&s.ID,
			&s.Name,
			&s.Description,
			&s.ImageURL,
			&s.CurrentEpisode,
			&s.TotalEpisodes,
		)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "could not read series data")
			return
		}
		seriesList = append(seriesList, s)
	}

	if err := rows.Err(); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "error while iterating series")
		return
	}

	writeJSON(w, http.StatusOK, seriesList)
}

// Handler GET /series/:id
func getSeriesByIDHandler(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromPath(r.URL.Path)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var s Series
	err = db.QueryRow(`
		SELECT id, name, description, image_url, current_episode, total_episodes
		FROM series2
		WHERE id = ?
	`, id).Scan(
		&s.ID,
		&s.Name,
		&s.Description,
		&s.ImageURL,
		&s.CurrentEpisode,
		&s.TotalEpisodes,
	)

	if err == sql.ErrNoRows {
		writeJSONError(w, http.StatusNotFound, "series not found")
		return
	}
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "could not fetch series")
		return
	}

	writeJSON(w, http.StatusOK, s)
}

// Handler POST /series
func createSeriesHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var s Series
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid json")
		return
	}

	errors := validateSeriesInput(s)
	if len(errors) > 0 {
		writeJSONValidationErrors(w, errors)
		return
	}

	result, err := db.Exec(`
		INSERT INTO series2 (name, description, image_url, total_episodes)
		VALUES (?, ?, ?, ?)
	`, s.Name, s.Description, s.ImageURL, s.TotalEpisodes)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "could not create series")
		return
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "could not retrieve inserted id")
		return
	}

	s.ID = int(lastID)
	s.CurrentEpisode = 1

	writeJSON(w, http.StatusCreated, s)
}

// Handler PUT /series/:id
func updateSeriesHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	id, err := getIDFromPath(r.URL.Path)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var s Series
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid json")
		return
	}

	errors := validateSeriesInput(s)
	if len(errors) > 0 {
		writeJSONValidationErrors(w, errors)
		return
	}

	result, err := db.Exec(`
		UPDATE series2
		SET name = ?, description = ?, image_url = ?, current_episode = ?, total_episodes = ?
		WHERE id = ?
	`, s.Name, s.Description, s.ImageURL, s.CurrentEpisode, s.TotalEpisodes, id)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "could not update series")
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "could not verify updated rows")
		return
	}

	if rowsAffected == 0 {
		writeJSONError(w, http.StatusNotFound, "series not found")
		return
	}

	s.ID = id
	writeJSON(w, http.StatusOK, s)
}

// HandlerDELETE /series/:id
func deleteSeriesHandler(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromPath(r.URL.Path)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid id")
		return
	}

	result, err := db.Exec(`DELETE FROM series2 WHERE id = ?`, id)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "could not delete series")
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "could not verify deleted rows")
		return
	}

	if rowsAffected == 0 {
		writeJSONError(w, http.StatusNotFound, "series not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
