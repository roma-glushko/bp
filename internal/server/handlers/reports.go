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

	sectionsParam := r.URL.Query().Get("sections")
	sec := report.Sections{
		Summary:       true,
		MonthAverages: true,
		WeekAverages:  true,
		DayAverages:   true,
		Sessions:      true,
		Readings:      true,
		Notes:         true,
	}

	if sectionsParam != "" {
		sec = report.Sections{}
		for _, s := range strings.Split(sectionsParam, ",") {
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
	}

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

	rpt := report.Generate(
		sessions,
		patient,
		device,
		from.Format("2006-01-02"),
		to.Format("2006-01-02"),
		time.Now().Format("2006-01-02"),
		sec,
	)

	writeJSON(w, http.StatusOK, rpt)
}
