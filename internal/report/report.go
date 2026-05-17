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

package report

import (
	"github.com/roma-glushko/bp/internal/domain"
)

type Report struct {
	From      string `json:"from"`
	To        string `json:"to"`
	Patient   string `json:"patient,omitempty"`
	Device    string `json:"device,omitempty"`
	Generated string `json:"generated"`

	Summary       *Summary       `json:"summary,omitempty"`
	MonthAverages []MonthAverage `json:"month_averages,omitempty"`
	WeekAverages  []WeekAverage  `json:"week_averages,omitempty"`
	DayAverages   []DayAverage   `json:"day_averages,omitempty"`
	Sessions      []SessionEntry `json:"sessions,omitempty"`
	Readings      []ReadingEntry `json:"readings,omitempty"`
}

type ReadingEntry struct {
	Date      string `json:"date"`
	Time      string `json:"time"`
	Period    string `json:"period"`
	Attempt   int    `json:"attempt"`
	Systolic  int    `json:"systolic"`
	Diastolic int    `json:"diastolic"`
	Pulse     *int   `json:"pulse,omitempty"`
}

type Sections struct {
	Summary       bool
	MonthAverages bool
	WeekAverages  bool
	DayAverages   bool
	Sessions      bool
	Readings      bool
	Notes         bool
}

func buildSessionEntries(sessions []domain.MeasurementSession) []SessionEntry {
	entries := make([]SessionEntry, 0, len(sessions))
	for _, s := range sessions {
		entries = append(entries, SessionEntry{
			MeasurementSession: s,
			Average:            domain.ComputeSessionAverage(s.Readings),
		})
	}
	return entries
}

func buildReadingEntries(sessions []domain.MeasurementSession) []ReadingEntry {
	var entries []ReadingEntry
	for _, s := range sessions {
		for _, rd := range s.Readings {
			entries = append(entries, ReadingEntry{
				Date:      s.MeasuredAt.Format("2006-01-02"),
				Time:      s.MeasuredAt.Format("15:04"),
				Period:    string(s.Period),
				Attempt:   rd.ReadingNo,
				Systolic:  rd.Systolic,
				Diastolic: rd.Diastolic,
				Pulse:     rd.Pulse,
			})
		}
	}
	return entries
}

func applyDailyNotes(dayAvgs []DayAverage, annotations *domain.NotesFile) {
	if annotations == nil {
		return
	}
	noteMap := make(map[string]string, len(annotations.DailyNotes))
	for _, dn := range annotations.DailyNotes {
		noteMap[dn.Date] = dn.Notes
	}
	for i := range dayAvgs {
		if n, ok := noteMap[dayAvgs[i].Date]; ok {
			dayAvgs[i].DayNote = n
		}
	}
}

func applyWeeklyNotes(weekAvgs []WeekAverage, annotations *domain.NotesFile) {
	if annotations == nil {
		return
	}
	noteMap := make(map[string]string, len(annotations.WeeklyNotes))
	for _, wn := range annotations.WeeklyNotes {
		noteMap[wn.Week] = wn.Notes
	}
	for i := range weekAvgs {
		if n, ok := noteMap[weekAvgs[i].Week]; ok {
			weekAvgs[i].Notes = n
		}
	}
}

func stripNotes(dayAvgs []DayAverage) {
	for i := range dayAvgs {
		dayAvgs[i].Notes = nil
		dayAvgs[i].DayNote = ""
	}
}

func Generate(sessions []domain.MeasurementSession, annotations *domain.NotesFile, patient, device, from, to, generated string, sec Sections) Report {
	r := Report{
		From:      from,
		To:        to,
		Patient:   patient,
		Device:    device,
		Generated: generated,
	}

	dayAvgs := ComputeDayAverages(sessions)

	if sec.Notes {
		applyDailyNotes(dayAvgs, annotations)
	}

	if sec.Summary {
		s := ComputeSummary(sessions)
		r.Summary = &s
	}

	if sec.MonthAverages {
		r.MonthAverages = ComputeMonthAverages(dayAvgs)
	}

	if sec.WeekAverages {
		r.WeekAverages = ComputeWeekAverages(dayAvgs)
		if sec.Notes {
			applyWeeklyNotes(r.WeekAverages, annotations)
		}
	}

	if sec.DayAverages {
		r.DayAverages = dayAvgs
		if !sec.Notes {
			stripNotes(r.DayAverages)
		}
	}

	if sec.Sessions {
		r.Sessions = buildSessionEntries(sessions)
	}

	if sec.Readings {
		r.Readings = buildReadingEntries(sessions)
	}

	return r
}
