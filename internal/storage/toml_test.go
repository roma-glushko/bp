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

package storage

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/roma-glushko/bp/internal/domain"
)

func intPtr(v int) *int { return &v }

func newTestStore(t *testing.T) *TomlStore {
	t.Helper()
	store, err := NewTomlStore(t.TempDir())
	require.NoError(t, err)
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

	require.NoError(t, store.CreateSession(&session))
	assert.NotEmpty(t, session.ID)

	got, err := store.GetSession(session.ID)
	require.NoError(t, err)

	assert.Equal(t, session.ID, got.ID)
	assert.Equal(t, "test session", got.Notes)
	assert.Len(t, got.Readings, 2)
	assert.NotEmpty(t, got.Readings[0].ID)
	assert.Equal(t, 1, got.Readings[0].ReadingNo)
}

func TestGetSessionNotFound(t *testing.T) {
	store := newTestStore(t)

	_, err := store.GetSession("nonexistent-id")
	assert.ErrorIs(t, err, ErrSessionNotFound)
}

func TestUpdateSession(t *testing.T) {
	store := newTestStore(t)
	session := makeSession(time.Date(2026, 5, 7, 9, 0, 0, 0, time.UTC))

	require.NoError(t, store.CreateSession(&session))

	session.Notes = "updated notes"
	session.Readings = append(session.Readings, domain.Reading{
		Systolic: 115, Diastolic: 78, Pulse: intPtr(72),
	})

	require.NoError(t, store.UpdateSession(&session))

	got, err := store.GetSession(session.ID)
	require.NoError(t, err)

	assert.Equal(t, "updated notes", got.Notes)
	assert.Len(t, got.Readings, 3)
}

func TestUpdateSessionCrossMonth(t *testing.T) {
	store := newTestStore(t)
	session := makeSession(time.Date(2026, 5, 31, 22, 0, 0, 0, time.UTC))

	require.NoError(t, store.CreateSession(&session))

	session.MeasuredAt = time.Date(2026, 6, 1, 8, 0, 0, 0, time.UTC)
	require.NoError(t, store.UpdateSession(&session))

	got, err := store.GetSession(session.ID)
	require.NoError(t, err)
	assert.True(t, got.MeasuredAt.Equal(session.MeasuredAt))

	maySessions, err := store.ListSessions(
		time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 5, 31, 23, 59, 59, 0, time.UTC),
	)
	require.NoError(t, err)

	for _, s := range maySessions {
		assert.NotEqual(t, session.ID, s.ID, "session still found in May after cross-month update")
	}
}

func TestUpdateSessionNotFound(t *testing.T) {
	store := newTestStore(t)
	session := makeSession(time.Date(2026, 5, 7, 9, 0, 0, 0, time.UTC))
	session.ID = "nonexistent-id"

	assert.ErrorIs(t, store.UpdateSession(&session), ErrSessionNotFound)
}

func TestDeleteSession(t *testing.T) {
	store := newTestStore(t)
	session := makeSession(time.Date(2026, 5, 7, 9, 0, 0, 0, time.UTC))

	require.NoError(t, store.CreateSession(&session))
	require.NoError(t, store.DeleteSession(session.ID))

	_, err := store.GetSession(session.ID)
	assert.ErrorIs(t, err, ErrSessionNotFound)
}

func TestDeleteSessionNotFound(t *testing.T) {
	store := newTestStore(t)
	assert.ErrorIs(t, store.DeleteSession("nonexistent-id"), ErrSessionNotFound)
}

func TestListSessions(t *testing.T) {
	store := newTestStore(t)

	may7 := makeSession(time.Date(2026, 5, 7, 9, 0, 0, 0, time.UTC))
	may15 := makeSession(time.Date(2026, 5, 15, 22, 0, 0, 0, time.UTC))
	jun1 := makeSession(time.Date(2026, 6, 1, 8, 0, 0, 0, time.UTC))

	for _, s := range []*domain.MeasurementSession{&may7, &may15, &jun1} {
		require.NoError(t, store.CreateSession(s))
	}

	// List all
	all, err := store.ListSessions(
		time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 6, 30, 23, 59, 59, 0, time.UTC),
	)
	require.NoError(t, err)
	assert.Len(t, all, 3)

	if assert.Len(t, all, 3) {
		assert.True(t, all[0].MeasuredAt.After(all[1].MeasuredAt), "not sorted descending")
	}

	// List May only
	mayOnly, err := store.ListSessions(
		time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2026, 5, 31, 23, 59, 59, 0, time.UTC),
	)
	require.NoError(t, err)
	assert.Len(t, mayOnly, 2)

	// List with no matches
	empty, err := store.ListSessions(
		time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 1, 31, 23, 59, 59, 0, time.UTC),
	)
	require.NoError(t, err)
	assert.Empty(t, empty)
}

func TestSettingsRoundTrip(t *testing.T) {
	store := newTestStore(t)

	settings, err := store.GetSettings()
	require.NoError(t, err)
	assert.Empty(t, settings.PatientName)

	settings.PatientName = "Roman Hlushko"
	settings.DefaultArm = domain.ArmLeft
	settings.DefaultPosition = domain.PositionSitting

	require.NoError(t, store.SaveSettings(settings))

	got, err := store.GetSettings()
	require.NoError(t, err)

	assert.Equal(t, "Roman Hlushko", got.PatientName)
	assert.Equal(t, domain.ArmLeft, got.DefaultArm)
	assert.False(t, got.UpdatedAt.IsZero())
}
