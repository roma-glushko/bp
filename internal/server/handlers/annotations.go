// Copyright 2025 Roma Hlushko
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package handlers

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/roma-glushko/bp/internal/domain"
	"github.com/roma-glushko/bp/internal/storage"
)

type AnnotationHandler struct {
	Store storage.Store
}

func (h *AnnotationHandler) List(w http.ResponseWriter, r *http.Request) {
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

	nf, err := h.Store.ListAnnotations(from, to)
	if err != nil {
		slog.Error("listing annotations", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to list annotations")
		return
	}

	writeJSON(w, http.StatusOK, nf)
}

func (h *AnnotationHandler) UpsertDaily(w http.ResponseWriter, r *http.Request) {
	date := r.PathValue("date")

	var body struct {
		Notes string `json:"notes"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if body.Notes == "" {
		if err := h.Store.DeleteDailyNote(date); err != nil {
			slog.Error("deleting daily note", "error", err)
			writeError(w, http.StatusInternalServerError, "failed to delete daily note")
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}

	note := &domain.DailyNote{
		Date:  date,
		Notes: body.Notes,
	}

	if err := h.Store.UpsertDailyNote(note); err != nil {
		slog.Error("upserting daily note", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to save daily note")
		return
	}

	writeJSON(w, http.StatusOK, note)
}

func (h *AnnotationHandler) DeleteDaily(w http.ResponseWriter, r *http.Request) {
	date := r.PathValue("date")

	if err := h.Store.DeleteDailyNote(date); err != nil {
		slog.Error("deleting daily note", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to delete daily note")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AnnotationHandler) UpsertWeekly(w http.ResponseWriter, r *http.Request) {
	week := r.PathValue("week")

	var body struct {
		Notes string `json:"notes"`
	}
	if err := readJSON(r, &body); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if body.Notes == "" {
		if err := h.Store.DeleteWeeklyNote(week); err != nil {
			slog.Error("deleting weekly note", "error", err)
			writeError(w, http.StatusInternalServerError, "failed to delete weekly note")
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}

	note := &domain.WeeklyNote{
		Week:  week,
		Notes: body.Notes,
	}

	if err := h.Store.UpsertWeeklyNote(note); err != nil {
		slog.Error("upserting weekly note", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to save weekly note")
		return
	}

	writeJSON(w, http.StatusOK, note)
}

func (h *AnnotationHandler) DeleteWeekly(w http.ResponseWriter, r *http.Request) {
	week := r.PathValue("week")

	if err := h.Store.DeleteWeeklyNote(week); err != nil {
		slog.Error("deleting weekly note", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to delete weekly note")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
