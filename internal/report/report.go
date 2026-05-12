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

func Generate(sessions []domain.MeasurementSession, patient, device, from, to, generated string, sec Sections) Report {
	r := Report{
		From:      from,
		To:        to,
		Patient:   patient,
		Device:    device,
		Generated: generated,
	}

	dayAvgs := ComputeDayAverages(sessions)

	if sec.Summary {
		s := ComputeSummary(sessions)
		r.Summary = &s
	}

	if sec.MonthAverages {
		r.MonthAverages = ComputeMonthAverages(dayAvgs)
	}

	if sec.WeekAverages {
		r.WeekAverages = ComputeWeekAverages(dayAvgs)
	}

	if sec.DayAverages {
		r.DayAverages = dayAvgs
		if !sec.Notes {
			for i := range r.DayAverages {
				r.DayAverages[i].Notes = nil
			}
		}
	}

	if sec.Sessions {
		entries := make([]SessionEntry, 0, len(sessions))
		for _, s := range sessions {
			entries = append(entries, SessionEntry{
				MeasurementSession: s,
				Average:            domain.ComputeSessionAverage(s.Readings),
			})
		}
		r.Sessions = entries
	}

	if sec.Readings {
		for _, s := range sessions {
			for _, rd := range s.Readings {
				r.Readings = append(r.Readings, ReadingEntry{
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
	}

	return r
}
