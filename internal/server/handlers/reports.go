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
	"strings"
	"time"

	"github.com/roma-glushko/bp/internal/report"
	"github.com/roma-glushko/bp/internal/storage"
)

type ReportHandler struct {
	Store storage.Store
}

func parseSections(param string) report.Sections {
	if param == "" {
		return report.Sections{
			Summary:       true,
			MonthAverages: true,
			WeekAverages:  true,
			DayAverages:   true,
			Sessions:      true,
			Readings:      true,
			Notes:         true,
		}
	}

	sec := report.Sections{}

	for _, s := range strings.Split(param, ",") {
		switch strings.TrimSpace(s) {
		case "summary":
			sec.Summary = true
		case "monthly":
			sec.MonthAverages = true
		case "weekly":
			sec.WeekAverages = true
		case "daily":
			sec.DayAverages = true
		case "sessions":
			sec.Sessions = true
		case "readings":
			sec.Readings = true
		case "notes":
			sec.Notes = true
		}
	}

	return sec
}

func (h *ReportHandler) Preview(w http.ResponseWriter, r *http.Request) {
	to := time.Now().UTC()
	from := to.AddDate(0, -1, 0)

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

	sec := parseSections(r.URL.Query().Get("sections"))

	sessions, err := h.Store.ListSessions(from, to)
	if err != nil {
		slog.Error("listing sessions for report", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to generate report")
		return
	}

	settings, _ := h.Store.GetSettings()
	patient := ""
	device := ""
	if settings != nil {
		patient = settings.PatientName
		device = settings.DeviceName
	}

	annotations, _ := h.Store.ListAnnotations(from, to)

	rpt := report.Generate(
		sessions,
		annotations,
		patient,
		device,
		from.Format("2006-01-02"),
		to.Format("2006-01-02"),
		time.Now().Format("2006-01-02"),
		sec,
	)

	writeJSON(w, http.StatusOK, rpt)
}
