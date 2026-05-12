package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/roma-glushko/bp/internal/domain"
	"github.com/roma-glushko/bp/internal/storage"
)

func newTestSettingsHandler(t *testing.T) *SettingsHandler {
	t.Helper()
	store, err := storage.NewTomlStore(t.TempDir())
	require.NoError(t, err)
	return &SettingsHandler{Store: store}
}

func TestGetDefaultSettings(t *testing.T) {
	h := newTestSettingsHandler(t)

	req := httptest.NewRequest(http.MethodGet, "/api/settings", nil)
	w := httptest.NewRecorder()
	h.Get(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var settings domain.Settings
	require.NoError(t, json.NewDecoder(w.Body).Decode(&settings))
	assert.Empty(t, settings.PatientName)
}

func TestUpdateAndGetSettings(t *testing.T) {
	h := newTestSettingsHandler(t)

	body := `{"patient_name": "Roman Hlushko", "default_arm": "left", "default_position": "sitting"}`
	req := httptest.NewRequest(http.MethodPost, "/api/settings", strings.NewReader(body))
	w := httptest.NewRecorder()
	h.Update(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	// Read back
	req = httptest.NewRequest(http.MethodGet, "/api/settings", nil)
	w = httptest.NewRecorder()
	h.Get(w, req)

	var got domain.Settings

	require.NoError(t, json.NewDecoder(w.Body).Decode(&got))
	assert.Equal(t, "Roman Hlushko", got.PatientName)
	assert.Equal(t, domain.ArmLeft, got.DefaultArm)
}
