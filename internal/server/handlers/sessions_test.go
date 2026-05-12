package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/roma-glushko/bp/internal/domain"
	"github.com/roma-glushko/bp/internal/storage"
)

func newTestHandler(t *testing.T) *SessionHandler {
	t.Helper()
	store, err := storage.NewTomlStore(t.TempDir())
	require.NoError(t, err)
	return &SessionHandler{Store: store}
}

const validSessionJSON = `{
	"measured_at": "2026-05-07T09:00:00Z",
	"period": "morning",
	"arm": "left",
	"position": "sitting",
	"readings": [
		{"systolic": 126, "diastolic": 86, "pulse": 78},
		{"systolic": 119, "diastolic": 83, "pulse": 81}
	],
	"notes": "test session"
}`

func createSession(t *testing.T, h *SessionHandler) sessionResponse {
	t.Helper()
	req := httptest.NewRequest("POST", "/api/sessions", strings.NewReader(validSessionJSON))
	w := httptest.NewRecorder()
	h.Create(w, req)
	require.Equal(t, http.StatusCreated, w.Code)

	var resp sessionResponse
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	return resp
}

func TestCreateSession(t *testing.T) {
	h := newTestHandler(t)
	resp := createSession(t, h)

	assert.NotEmpty(t, resp.ID)
	assert.Len(t, resp.Readings, 2)
	assert.True(t, resp.Average.AvgSystolic.Equal(decimal.NewFromFloat(122.5)))
	assert.Equal(t, domain.BPStatusElevated, resp.Average.Indicator.Status)
}

func TestCreateSessionInvalid(t *testing.T) {
	h := newTestHandler(t)

	body := `{"period": "morning", "readings": []}`
	req := httptest.NewRequest("POST", "/api/sessions", strings.NewReader(body))
	w := httptest.NewRecorder()
	h.Create(w, req)

	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

	var resp map[string]any
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	assert.Contains(t, resp, "errors")
}

func TestCreateSessionBadJSON(t *testing.T) {
	h := newTestHandler(t)

	req := httptest.NewRequest("POST", "/api/sessions", strings.NewReader("not json"))
	w := httptest.NewRecorder()
	h.Create(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestListSessions(t *testing.T) {
	h := newTestHandler(t)
	createSession(t, h)

	req := httptest.NewRequest("GET", "/api/sessions?from=2026-01-01T00:00:00Z&to=2026-12-31T23:59:59Z", nil)
	w := httptest.NewRecorder()
	h.List(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp listResponse
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	assert.Len(t, resp.Sessions, 1)
}

func TestListSessionsDefaultRange(t *testing.T) {
	h := newTestHandler(t)

	req := httptest.NewRequest("GET", "/api/sessions", nil)
	w := httptest.NewRecorder()
	h.List(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetSession(t *testing.T) {
	h := newTestHandler(t)
	created := createSession(t, h)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/sessions/{id}", h.Get)

	req := httptest.NewRequest("GET", "/api/sessions/"+created.ID, nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var got sessionResponse
	require.NoError(t, json.NewDecoder(w.Body).Decode(&got))
	assert.Equal(t, created.ID, got.ID)
}

func TestGetSessionNotFound(t *testing.T) {
	h := newTestHandler(t)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/sessions/{id}", h.Get)

	req := httptest.NewRequest("GET", "/api/sessions/nonexistent", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestUpdateSession(t *testing.T) {
	h := newTestHandler(t)
	created := createSession(t, h)

	updatedJSON := `{
		"measured_at": "2026-05-07T09:00:00Z",
		"period": "morning",
		"arm": "right",
		"position": "sitting",
		"readings": [{"systolic": 115, "diastolic": 75, "pulse": 70}],
		"notes": "updated"
	}`

	mux := http.NewServeMux()
	mux.HandleFunc("PUT /api/sessions/{id}", h.Update)

	req := httptest.NewRequest("PUT", "/api/sessions/"+created.ID, strings.NewReader(updatedJSON))
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp sessionResponse
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	assert.Equal(t, "updated", resp.Notes)
	assert.Len(t, resp.Readings, 1)
	assert.False(t, resp.CreatedAt.IsZero(), "CreatedAt should be preserved from original")
}

func TestDeleteSession(t *testing.T) {
	h := newTestHandler(t)
	created := createSession(t, h)

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /api/sessions/{id}", h.Delete)
	mux.HandleFunc("GET /api/sessions/{id}", h.Get)

	req := httptest.NewRequest("DELETE", "/api/sessions/"+created.ID, nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)

	// Verify gone
	req = httptest.NewRequest("GET", "/api/sessions/"+created.ID, nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteSessionNotFound(t *testing.T) {
	h := newTestHandler(t)

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /api/sessions/{id}", h.Delete)

	req := httptest.NewRequest("DELETE", "/api/sessions/nonexistent", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
