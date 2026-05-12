package storage

import (
	"errors"
	"testing"
	"time"

	"github.com/roma-glushko/bp/internal/domain"
)

func intPtr(v int) *int { return &v }

func newTestStore(t *testing.T) *TomlStore {
	t.Helper()
	store, err := NewTomlStore(t.TempDir())
	if err != nil {
		t.Fatalf("NewTomlStore() error: %v", err)
	}
	return store
}

func makeSession(measuredAt time.Time) domain.MeasurementSession {
	return domain.MeasurementSession{
		MeasuredAt: measuredAt,
		Period:     domain.PeriodMorning,
		Arm:        domain.ArmLeft,
		Position:   domain.PositionSitting,
		Readings: []domain.Reading{
			{Systolic: 126, Diastolic: 86, Pulse: intPtr(78)},
			{Systolic: 119, Diastolic: 83, Pulse: intPtr(81)},
		},
		Notes: "test session",
	}
}

func TestCreateAndGetSession(t *testing.T) {
	store := newTestStore(t)
	session := makeSession(time.Date(2026, 5, 7, 9, 0, 0, 0, time.UTC))

	if err := store.CreateSession(&session); err != nil {
		t.Fatalf("CreateSession() error: %v", err)
	}

	if session.ID == "" {
		t.Fatal("session ID not set after create")
	}

	got, err := store.GetSession(session.ID)
	if err != nil {
		t.Fatalf("GetSession() error: %v", err)
	}

	if got.ID != session.ID {
		t.Errorf("ID = %q, want %q", got.ID, session.ID)
	}
	if got.Notes != "test session" {
		t.Errorf("Notes = %q, want %q", got.Notes, "test session")
	}
	if len(got.Readings) != 2 {
		t.Errorf("Readings count = %d, want 2", len(got.Readings))
	}
	if got.Readings[0].ID == "" {
		t.Error("reading ID not set after create")
	}
	if got.Readings[0].ReadingNo != 1 {
		t.Errorf("reading_no = %d, want 1", got.Readings[0].ReadingNo)
	}
}

func TestGetSessionNotFound(t *testing.T) {
	store := newTestStore(t)

	_, err := store.GetSession("nonexistent-id")
	if !errors.Is(err, ErrSessionNotFound) {
		t.Errorf("GetSession() error = %v, want ErrSessionNotFound", err)
	}
}

func TestUpdateSession(t *testing.T) {
	store := newTestStore(t)
	session := makeSession(time.Date(2026, 5, 7, 9, 0, 0, 0, time.UTC))

	if err := store.CreateSession(&session); err != nil {
		t.Fatalf("CreateSession() error: %v", err)
	}

	session.Notes = "updated notes"
	session.Readings = append(session.Readings, domain.Reading{
		Systolic: 115, Diastolic: 78, Pulse: intPtr(72),
	})

	if err := store.UpdateSession(&session); err != nil {
		t.Fatalf("UpdateSession() error: %v", err)
	}

	got, err := store.GetSession(session.ID)
	if err != nil {
		t.Fatalf("GetSession() error: %v", err)
	}

	if got.Notes != "updated notes" {
		t.Errorf("Notes = %q, want %q", got.Notes, "updated notes")
	}
	if len(got.Readings) != 3 {
		t.Errorf("Readings count = %d, want 3", len(got.Readings))
	}
}

func TestUpdateSessionCrossMonth(t *testing.T) {
	store := newTestStore(t)
	session := makeSession(time.Date(2026, 5, 31, 22, 0, 0, 0, time.UTC))

	if err := store.CreateSession(&session); err != nil {
		t.Fatalf("CreateSession() error: %v", err)
	}

	// Move to June
	session.MeasuredAt = time.Date(2026, 6, 1, 8, 0, 0, 0, time.UTC)
	if err := store.UpdateSession(&session); err != nil {
		t.Fatalf("UpdateSession() error: %v", err)
	}

	// Should be findable
	got, err := store.GetSession(session.ID)
	if err != nil {
		t.Fatalf("GetSession() error: %v", err)
	}
	if !got.MeasuredAt.Equal(session.MeasuredAt) {
		t.Errorf("MeasuredAt = %v, want %v", got.MeasuredAt, session.MeasuredAt)
	}

	// Should not appear in May listing
	maySessions, err := store.ListSessions(
		time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 5, 31, 23, 59, 59, 0, time.UTC),
	)
	if err != nil {
		t.Fatalf("ListSessions() error: %v", err)
	}
	for _, s := range maySessions {
		if s.ID == session.ID {
			t.Error("session still found in May after cross-month update")
		}
	}
}

func TestUpdateSessionNotFound(t *testing.T) {
	store := newTestStore(t)
	session := makeSession(time.Date(2026, 5, 7, 9, 0, 0, 0, time.UTC))
	session.ID = "nonexistent-id"

	err := store.UpdateSession(&session)
	if !errors.Is(err, ErrSessionNotFound) {
		t.Errorf("UpdateSession() error = %v, want ErrSessionNotFound", err)
	}
}

func TestDeleteSession(t *testing.T) {
	store := newTestStore(t)
	session := makeSession(time.Date(2026, 5, 7, 9, 0, 0, 0, time.UTC))

	if err := store.CreateSession(&session); err != nil {
		t.Fatalf("CreateSession() error: %v", err)
	}

	if err := store.DeleteSession(session.ID); err != nil {
		t.Fatalf("DeleteSession() error: %v", err)
	}

	_, err := store.GetSession(session.ID)
	if !errors.Is(err, ErrSessionNotFound) {
		t.Errorf("GetSession() after delete: error = %v, want ErrSessionNotFound", err)
	}
}

func TestDeleteSessionNotFound(t *testing.T) {
	store := newTestStore(t)
	err := store.DeleteSession("nonexistent-id")
	if !errors.Is(err, ErrSessionNotFound) {
		t.Errorf("DeleteSession() error = %v, want ErrSessionNotFound", err)
	}
}

func TestListSessions(t *testing.T) {
	store := newTestStore(t)

	may7 := makeSession(time.Date(2026, 5, 7, 9, 0, 0, 0, time.UTC))
	may15 := makeSession(time.Date(2026, 5, 15, 22, 0, 0, 0, time.UTC))
	jun1 := makeSession(time.Date(2026, 6, 1, 8, 0, 0, 0, time.UTC))

	for _, s := range []*domain.MeasurementSession{&may7, &may15, &jun1} {
		if err := store.CreateSession(s); err != nil {
			t.Fatalf("CreateSession() error: %v", err)
		}
	}

	// List all
	all, err := store.ListSessions(
		time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 6, 30, 23, 59, 59, 0, time.UTC),
	)
	if err != nil {
		t.Fatalf("ListSessions() error: %v", err)
	}
	if len(all) != 3 {
		t.Errorf("ListSessions(all) count = %d, want 3", len(all))
	}
	// Verify descending order
	if len(all) >= 2 && all[0].MeasuredAt.Before(all[1].MeasuredAt) {
		t.Error("ListSessions() not sorted descending by measured_at")
	}

	// List May only
	mayOnly, err := store.ListSessions(
		time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 5, 31, 23, 59, 59, 0, time.UTC),
	)
	if err != nil {
		t.Fatalf("ListSessions() error: %v", err)
	}
	if len(mayOnly) != 2 {
		t.Errorf("ListSessions(May) count = %d, want 2", len(mayOnly))
	}

	// List with no matches
	empty, err := store.ListSessions(
		time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 1, 31, 23, 59, 59, 0, time.UTC),
	)
	if err != nil {
		t.Fatalf("ListSessions() error: %v", err)
	}
	if len(empty) != 0 {
		t.Errorf("ListSessions(empty) count = %d, want 0", len(empty))
	}
}

func TestSettingsRoundTrip(t *testing.T) {
	store := newTestStore(t)

	// Get default settings (no file exists)
	settings, err := store.GetSettings()
	if err != nil {
		t.Fatalf("GetSettings() error: %v", err)
	}
	if settings.PatientName != "" {
		t.Errorf("default PatientName = %q, want empty", settings.PatientName)
	}

	// Save and reload
	settings.PatientName = "Roman Hlushko"
	settings.DefaultArm = domain.ArmLeft
	settings.DefaultPosition = domain.PositionSitting

	if err := store.SaveSettings(settings); err != nil {
		t.Fatalf("SaveSettings() error: %v", err)
	}

	got, err := store.GetSettings()
	if err != nil {
		t.Fatalf("GetSettings() error: %v", err)
	}
	if got.PatientName != "Roman Hlushko" {
		t.Errorf("PatientName = %q, want %q", got.PatientName, "Roman Hlushko")
	}
	if got.DefaultArm != domain.ArmLeft {
		t.Errorf("DefaultArm = %q, want %q", got.DefaultArm, domain.ArmLeft)
	}
	if got.UpdatedAt.IsZero() {
		t.Error("UpdatedAt not set after save")
	}
}
