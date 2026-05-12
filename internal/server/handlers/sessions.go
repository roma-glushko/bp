package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/roma-glushko/bp/internal/domain"
	"github.com/roma-glushko/bp/internal/storage"
)

type SessionHandler struct {
	Store storage.Store
}

type sessionResponse struct {
	domain.MeasurementSession
	Average domain.SessionAverage `json:"average"`
}

type listResponse struct {
	Sessions []sessionResponse `json:"sessions"`
}

func enrichSession(s domain.MeasurementSession) sessionResponse {
	return sessionResponse{
		MeasurementSession: s,
		Average:            domain.ComputeSessionAverage(s.Readings),
	}
}

func (h *SessionHandler) List(w http.ResponseWriter, r *http.Request) {
	to := time.Now().UTC()
	from := to.AddDate(0, 0, -30)

	if v := r.URL.Query().Get("from"); v != "" {
		t, err := time.Parse(time.RFC3339, v)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid 'from' date format, use RFC 3339")
			return
		}
		from = t
	}
	if v := r.URL.Query().Get("to"); v != "" {
		t, err := time.Parse(time.RFC3339, v)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid 'to' date format, use RFC 3339")
			return
		}
		to = t
	}

	sessions, err := h.Store.ListSessions(from, to)
	if err != nil {
		slog.Error("listing sessions", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to list sessions")
		return
	}

	resp := listResponse{Sessions: make([]sessionResponse, 0, len(sessions))}
	for _, s := range sessions {
		resp.Sessions = append(resp.Sessions, enrichSession(s))
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *SessionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var session domain.MeasurementSession
	if err := readJSON(r, &session); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if errs := domain.ValidateSession(session); len(errs) > 0 {
		writeValidationErrors(w, errs)
		return
	}

	if err := h.Store.CreateSession(&session); err != nil {
		slog.Error("creating session", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to create session")
		return
	}

	writeJSON(w, http.StatusCreated, enrichSession(session))
}

func (h *SessionHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	session, err := h.Store.GetSession(id)
	if err != nil {
		if errors.Is(err, storage.ErrSessionNotFound) {
			writeError(w, http.StatusNotFound, "session not found")
			return
		}
		slog.Error("getting session", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to get session")
		return
	}

	writeJSON(w, http.StatusOK, enrichSession(*session))
}

func (h *SessionHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var session domain.MeasurementSession
	if err := readJSON(r, &session); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	session.ID = id

	if errs := domain.ValidateSession(session); len(errs) > 0 {
		writeValidationErrors(w, errs)
		return
	}

	// Preserve original CreatedAt
	existing, err := h.Store.GetSession(id)
	if err != nil {
		if errors.Is(err, storage.ErrSessionNotFound) {
			writeError(w, http.StatusNotFound, "session not found")
			return
		}
		slog.Error("getting session for update", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to get session")
		return
	}
	session.CreatedAt = existing.CreatedAt

	if err := h.Store.UpdateSession(&session); err != nil {
		if errors.Is(err, storage.ErrSessionNotFound) {
			writeError(w, http.StatusNotFound, "session not found")
			return
		}

		slog.Error("updating session", "error", err)

		writeError(w, http.StatusInternalServerError, "failed to update session")

		return
	}

	writeJSON(w, http.StatusOK, enrichSession(session))
}

func (h *SessionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := h.Store.DeleteSession(id); err != nil {
		if errors.Is(err, storage.ErrSessionNotFound) {
			writeError(w, http.StatusNotFound, "session not found")
			return
		}

		slog.Error("deleting session", "error", err)

		writeError(w, http.StatusInternalServerError, "failed to delete session")

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
