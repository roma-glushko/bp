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
	"errors"
	"time"

	"github.com/roma-glushko/bp/internal/domain"
)

var ErrSessionNotFound = errors.New("session not found")

type Store interface {
	CreateSession(session *domain.MeasurementSession) error
	GetSession(id string) (*domain.MeasurementSession, error)
	UpdateSession(session *domain.MeasurementSession) error
	DeleteSession(id string) error
	ListSessions(from, to time.Time) ([]domain.MeasurementSession, error)

	GetSettings() (*domain.Settings, error)
	SaveSettings(settings *domain.Settings) error

	ListAnnotations(from, to time.Time) (*domain.NotesFile, error)
	UpsertDailyNote(note *domain.DailyNote) error
	DeleteDailyNote(date string) error
	UpsertWeeklyNote(note *domain.WeeklyNote) error
	DeleteWeeklyNote(week string) error
}
