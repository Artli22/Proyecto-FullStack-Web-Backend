package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

func getIDFromPath(path string) (int, error) {
	idStr := strings.TrimPrefix(path, "/series2/")
	return strconv.Atoi(idStr)
}

// Handlers get todas las series
func getSeriesHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
		SELECT id, name, description, image_url, current_episode, total_episodes
		FROM series2
	`)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
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
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		seriesList = append(seriesList, s)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(seriesList)
}

// Handler get series por id
func getSeriesByIDHandler(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
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
		http.Error(w, "series not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(s)
}

// Handler para crear una nueva serie
func createSeriesHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var s Series
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(s.Name) == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	if s.TotalEpisodes <= 0 {
		http.Error(w, "total_episodes must be greater than 0", http.StatusBadRequest)
		return
	}

	result, err := db.Exec(`
		INSERT INTO series2 (name, description, image_url, total_episodes)
		VALUES (?, ?, ?, ?)
	`, s.Name, s.Description, s.ImageURL, s.TotalEpisodes)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	s.ID = int(lastID)
	s.CurrentEpisode = 1

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(s)
}

// Handler para actualizar una serie ya creada
func updateSeriesHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	id, err := getIDFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var s Series
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(s.Name) == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	if s.TotalEpisodes <= 0 {
		http.Error(w, "total_episodes must be greater than 0", http.StatusBadRequest)
		return
	}

	result, err := db.Exec(`
		UPDATE series2
		SET name = ?, description = ?, image_url = ?, current_episode = ?, total_episodes = ?
		WHERE id = ?
	`, s.Name, s.Description, s.ImageURL, s.CurrentEpisode, s.TotalEpisodes, id)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "series not found", http.StatusNotFound)
		return
	}

	s.ID = id

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(s)
}

// Handler parar eliminar una serie ya creada
func deleteSeriesHandler(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	result, err := db.Exec(`DELETE FROM series2 WHERE id = ?`, id)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "series not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
