package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type enquiryRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

type statusRequest struct {
	Status string `json:"status"`
}

func registerRoutes(mux *http.ServeMux, db *sql.DB) {
	mux.HandleFunc("POST /api/enquiries", createEnquiry(db))
	mux.HandleFunc("GET /api/enquiries", getEnquiries(db))
	mux.HandleFunc("GET /api/enquiry/{id}", getEnquiry(db))
	mux.HandleFunc("PATCH /api/enquiries/{id}", patchEnquiry(db))
	mux.HandleFunc("POST /api/enquiries/{id}/backdate", backdateHandler(db))
}

func createEnquiry(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req enquiryRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		req.Name = strings.TrimSpace(req.Name)
		req.Email = strings.TrimSpace(req.Email)
		req.Type = strings.TrimSpace(req.Type)
		req.Message = strings.TrimSpace(req.Message)

		if req.Name == "" || req.Email == "" || req.Type == "" || req.Message == "" {
			http.Error(w, "all fields are required", http.StatusBadRequest)
			return
		}
		if req.Type != "Membership" && req.Type != "Society Day" && req.Type != "Event" {
			http.Error(w, "type must be Membership, Society Day, or Event", http.StatusBadRequest)
			return
		}

		id, err := insertEnquiry(db, req.Name, req.Email, req.Type, req.Message)
		if err != nil {
			http.Error(w, "could not save enquiry", http.StatusInternalServerError)
			return
		}
		logConfirmationEmail(req.Name, req.Email, req.Type)
		writeJSON(w, http.StatusCreated, map[string]int64{"id": id})
	}
}

func getEnquiries(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		list, err := listEnquiries(db)
		if err != nil {
			http.Error(w, "could not list enquiries", http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, list)
	}
}

func getEnquiry(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := parseID(w, r)
		if !ok {
			return
		}
		enquiry, err := listEnquiry(db, id)
		if err == sql.ErrNoRows {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, "could not get enquiry", http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, enquiry)
	}
}

func patchEnquiry(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := parseID(w, r)
		if !ok {
			return
		}
		var req statusRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		if req.Status != "Contacted" {
			http.Error(w, `status must be "Contacted"`, http.StatusBadRequest)
			return
		}
		found, err := updateStatus(db, id, req.Status)
		if err != nil {
			http.Error(w, "could not update enquiry", http.StatusInternalServerError)
			return
		}
		if !found {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
	}
}

func backdateHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := parseID(w, r)
		if !ok {
			return
		}
		found, err := backdateEnquiry(db, id)
		if err != nil {
			http.Error(w, "could not backdate enquiry", http.StatusInternalServerError)
			return
		}
		if !found {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
	}
}

func parseID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id < 1 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return 0, false
	}
	return id, true
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
