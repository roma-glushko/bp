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
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
	toml "github.com/pelletier/go-toml/v2"

	"github.com/roma-glushko/bp/internal/config"
	"github.com/roma-glushko/bp/internal/domain"
)

type TomlStore struct {
	dataDir string
	mu      sync.RWMutex
}

func NewTomlStore(dataDir string) (*TomlStore, error) {
	if err := config.EnsureDataDir(dataDir); err != nil {
		return nil, err
	}
	return &TomlStore{dataDir: dataDir}, nil
}

func (s *TomlStore) monthFilePath(t time.Time) string {
	return filepath.Join(config.SessionsDir(s.dataDir), config.MonthFileName(t))
}

func (s *TomlStore) loadMonthFile(path string) (*domain.MonthFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &domain.MonthFile{}, nil
		}
		return nil, fmt.Errorf("reading month file: %w", err)
	}

	var mf domain.MonthFile
	if err := toml.Unmarshal(data, &mf); err != nil {
		return nil, fmt.Errorf("parsing month file %s: %w", path, err)
	}
	return &mf, nil
}

func (s *TomlStore) saveMonthFile(path string, mf *domain.MonthFile) error {
	data, err := toml.Marshal(mf)
	if err != nil {
		return fmt.Errorf("marshaling month file: %w", err)
	}

	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0o640); err != nil {
		return fmt.Errorf("writing temp file: %w", err)
	}

	if err := os.Rename(tmpPath, path); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("renaming temp file: %w", err)
	}
	return nil
}

func (s *TomlStore) allMonthFiles() ([]string, error) {
	sessionsDir := config.SessionsDir(s.dataDir)
	entries, err := os.ReadDir(sessionsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("reading sessions directory: %w", err)
	}

	var paths []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".toml" {
			paths = append(paths, filepath.Join(sessionsDir, e.Name()))
		}
	}
	return paths, nil
}

func (s *TomlStore) CreateSession(session *domain.MeasurementSession) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UTC()
	session.ID = uuid.New().String()
	session.CreatedAt = now
	session.UpdatedAt = now

	for i := range session.Readings {
		session.Readings[i].ID = uuid.New().String()
		session.Readings[i].ReadingNo = i + 1
		session.Readings[i].CreatedAt = now
	}

	path := s.monthFilePath(session.MeasuredAt)
	mf, err := s.loadMonthFile(path)
	if err != nil {
		return err
	}

	mf.Sessions = append(mf.Sessions, *session)
	return s.saveMonthFile(path, mf)
}

func (s *TomlStore) GetSession(id string) (*domain.MeasurementSession, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	paths, err := s.allMonthFiles()
	if err != nil {
		return nil, err
	}

	for _, path := range paths {
		mf, err := s.loadMonthFile(path)
		if err != nil {
			return nil, err
		}

		for i := range mf.Sessions {
			if mf.Sessions[i].ID == id {
				return &mf.Sessions[i], nil
			}
		}
	}

	return nil, ErrSessionNotFound
}

func (s *TomlStore) UpdateSession(session *domain.MeasurementSession) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	session.UpdatedAt = time.Now().UTC()

	// Find and remove the old session
	paths, err := s.allMonthFiles()
	if err != nil {
		return err
	}

	var oldPath string
	for _, path := range paths {
		mf, err := s.loadMonthFile(path)
		if err != nil {
			return err
		}

		for i := range mf.Sessions {
			if mf.Sessions[i].ID == session.ID {
				oldPath = path
				mf.Sessions = append(mf.Sessions[:i], mf.Sessions[i+1:]...)

				if err := s.saveMonthFile(path, mf); err != nil {
					return err
				}

				break
			}
		}
		if oldPath != "" {
			break
		}
	}

	if oldPath == "" {
		return ErrSessionNotFound
	}

	// Insert into the (possibly different) month file
	newPath := s.monthFilePath(session.MeasuredAt)
	mf, err := s.loadMonthFile(newPath)
	if err != nil {
		return err
	}
	mf.Sessions = append(mf.Sessions, *session)
	return s.saveMonthFile(newPath, mf)
}

func (s *TomlStore) DeleteSession(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	paths, err := s.allMonthFiles()
	if err != nil {
		return err
	}

	for _, path := range paths {
		mf, err := s.loadMonthFile(path)
		if err != nil {
			return err
		}

		for i := range mf.Sessions {
			if mf.Sessions[i].ID == id {
				mf.Sessions = append(mf.Sessions[:i], mf.Sessions[i+1:]...)
				return s.saveMonthFile(path, mf)
			}
		}
	}

	return ErrSessionNotFound
}

func (s *TomlStore) ListSessions(from, to time.Time) ([]domain.MeasurementSession, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	paths, err := s.allMonthFiles()
	if err != nil {
		return nil, err
	}

	var result []domain.MeasurementSession

	for _, path := range paths {
		// Quick check: does this month file potentially overlap the range?
		base := filepath.Base(path)
		monthStr := base[:len(base)-len(".toml")]
		fileMonth, err := time.Parse("2006-01", monthStr)
		if err != nil {
			continue
		}
		fileMonthEnd := fileMonth.AddDate(0, 1, 0)
		if fileMonth.After(to) || fileMonthEnd.Before(from) {
			continue
		}

		mf, err := s.loadMonthFile(path)
		if err != nil {
			return nil, err
		}

		for _, sess := range mf.Sessions {
			if !sess.MeasuredAt.Before(from) && !sess.MeasuredAt.After(to) {
				result = append(result, sess)
			}
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].MeasuredAt.After(result[j].MeasuredAt)
	})

	return result, nil
}

func (s *TomlStore) loadNotesFile() (*domain.NotesFile, error) {
	path := config.NotesPath(s.dataDir)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &domain.NotesFile{}, nil
		}
		return nil, fmt.Errorf("reading notes file: %w", err)
	}

	var nf domain.NotesFile
	if err := toml.Unmarshal(data, &nf); err != nil {
		return nil, fmt.Errorf("parsing notes file: %w", err)
	}
	return &nf, nil
}

func (s *TomlStore) saveNotesFile(nf *domain.NotesFile) error {
	path := config.NotesPath(s.dataDir)
	data, err := toml.Marshal(nf)
	if err != nil {
		return fmt.Errorf("marshaling notes file: %w", err)
	}

	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0o640); err != nil {
		return fmt.Errorf("writing temp notes file: %w", err)
	}

	if err := os.Rename(tmpPath, path); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("renaming temp notes file: %w", err)
	}
	return nil
}

func isoWeekRange(week string) (start, end time.Time, err error) {
	var year, w int
	if _, err := fmt.Sscanf(week, "%d-W%d", &year, &w); err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid week format: %s", week)
	}

	jan4 := time.Date(year, 1, 4, 0, 0, 0, 0, time.UTC)
	_, jan4Week := jan4.ISOWeek()
	start = jan4.AddDate(0, 0, (w-jan4Week)*7-int(jan4.Weekday()-time.Monday))
	if jan4.Weekday() == time.Sunday {
		start = start.AddDate(0, 0, -6)
	}
	end = start.AddDate(0, 0, 6)
	return start, end, nil
}

func (s *TomlStore) ListAnnotations(from, to time.Time) (*domain.NotesFile, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	nf, err := s.loadNotesFile()
	if err != nil {
		return nil, err
	}

	fromDate := from.Format("2006-01-02")
	toDate := to.Format("2006-01-02")

	var dailyNotes []domain.DailyNote
	for _, dn := range nf.DailyNotes {
		if dn.Date >= fromDate && dn.Date <= toDate {
			dailyNotes = append(dailyNotes, dn)
		}
	}

	var weeklyNotes []domain.WeeklyNote
	for _, wn := range nf.WeeklyNotes {
		wStart, wEnd, err := isoWeekRange(wn.Week)
		if err != nil {
			continue
		}
		if !wEnd.Before(from) && !wStart.After(to) {
			weeklyNotes = append(weeklyNotes, wn)
		}
	}

	return &domain.NotesFile{
		DailyNotes:  dailyNotes,
		WeeklyNotes: weeklyNotes,
	}, nil
}

func (s *TomlStore) UpsertDailyNote(note *domain.DailyNote) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	nf, err := s.loadNotesFile()
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	found := false
	for i := range nf.DailyNotes {
		if nf.DailyNotes[i].Date == note.Date {
			nf.DailyNotes[i].Notes = note.Notes
			nf.DailyNotes[i].UpdatedAt = now
			found = true
			break
		}
	}

	if !found {
		note.CreatedAt = now
		note.UpdatedAt = now
		nf.DailyNotes = append(nf.DailyNotes, *note)
	}

	return s.saveNotesFile(nf)
}

func (s *TomlStore) DeleteDailyNote(date string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	nf, err := s.loadNotesFile()
	if err != nil {
		return err
	}

	for i := range nf.DailyNotes {
		if nf.DailyNotes[i].Date == date {
			nf.DailyNotes = append(nf.DailyNotes[:i], nf.DailyNotes[i+1:]...)
			return s.saveNotesFile(nf)
		}
	}

	return nil
}

func (s *TomlStore) UpsertWeeklyNote(note *domain.WeeklyNote) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	nf, err := s.loadNotesFile()
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	found := false
	for i := range nf.WeeklyNotes {
		if nf.WeeklyNotes[i].Week == note.Week {
			nf.WeeklyNotes[i].Notes = note.Notes
			nf.WeeklyNotes[i].UpdatedAt = now
			found = true
			break
		}
	}

	if !found {
		note.CreatedAt = now
		note.UpdatedAt = now
		nf.WeeklyNotes = append(nf.WeeklyNotes, *note)
	}

	return s.saveNotesFile(nf)
}

func (s *TomlStore) DeleteWeeklyNote(week string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	nf, err := s.loadNotesFile()
	if err != nil {
		return err
	}

	for i := range nf.WeeklyNotes {
		if nf.WeeklyNotes[i].Week == week {
			nf.WeeklyNotes = append(nf.WeeklyNotes[:i], nf.WeeklyNotes[i+1:]...)
			return s.saveNotesFile(nf)
		}
	}

	return nil
}

func (s *TomlStore) GetSettings() (*domain.Settings, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	path := config.SettingsPath(s.dataDir)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &domain.Settings{}, nil
		}
		return nil, fmt.Errorf("reading settings: %w", err)
	}

	var sf domain.SettingsFile
	if err := toml.Unmarshal(data, &sf); err != nil {
		return nil, fmt.Errorf("parsing settings: %w", err)
	}
	return &sf.Settings, nil
}

func (s *TomlStore) SaveSettings(settings *domain.Settings) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	settings.UpdatedAt = time.Now().UTC()

	sf := domain.SettingsFile{Settings: *settings}
	data, err := toml.Marshal(&sf)
	if err != nil {
		return fmt.Errorf("marshaling settings: %w", err)
	}

	path := config.SettingsPath(s.dataDir)
	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0o640); err != nil {
		return fmt.Errorf("writing temp settings file: %w", err)
	}

	if err := os.Rename(tmpPath, path); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("renaming temp settings file: %w", err)
	}
	return nil
}
