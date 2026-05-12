package handlers

import (
	"log/slog"
	"net/http"

	"github.com/roma-glushko/bp/internal/domain"
	"github.com/roma-glushko/bp/internal/storage"
)

type SettingsHandler struct {
	Store storage.Store
}

func (h *SettingsHandler) Get(w http.ResponseWriter, _ *http.Request) {
	settings, err := h.Store.GetSettings()
	if err != nil {
		slog.Error("getting settings", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to get settings")
		return
	}

	writeJSON(w, http.StatusOK, settings)
}

func (h *SettingsHandler) Update(w http.ResponseWriter, r *http.Request) {
	var settings domain.Settings

	if err := readJSON(r, &settings); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.Store.SaveSettings(&settings); err != nil {
		slog.Error("saving settings", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to save settings")
		return
	}

	writeJSON(w, http.StatusOK, settings)
}
