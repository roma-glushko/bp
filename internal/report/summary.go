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
	"github.com/shopspring/decimal"

	"github.com/roma-glushko/bp/internal/domain"
)

type Summary struct {
	DaysMeasured int           `json:"days_measured"`
	SessionCount int           `json:"session_count"`
	ReadingCount int           `json:"reading_count"`
	OverallAvg   PeriodAverage `json:"overall_avg"`
	MorningAvg   PeriodAverage `json:"morning_avg"`
	EveningAvg   PeriodAverage `json:"evening_avg"`
	HighestSys   int           `json:"highest_systolic"`
	LowestSys    int           `json:"lowest_systolic"`
	HighestDia   int           `json:"highest_diastolic"`
	LowestDia    int           `json:"lowest_diastolic"`
	HighestPulse *int          `json:"highest_pulse,omitempty"`
	LowestPulse  *int          `json:"lowest_pulse,omitempty"`
}

type readingExtremes struct {
	count               int
	highSys, lowSys     int
	highDia, lowDia     int
	highPulse, lowPulse *int
}

func (ext *readingExtremes) trackReading(r domain.Reading) {
	ext.count++

	if r.Systolic > ext.highSys {
		ext.highSys = r.Systolic
	}
	if r.Systolic < ext.lowSys {
		ext.lowSys = r.Systolic
	}
	if r.Diastolic > ext.highDia {
		ext.highDia = r.Diastolic
	}
	if r.Diastolic < ext.lowDia {
		ext.lowDia = r.Diastolic
	}

	if r.Pulse != nil {
		if ext.highPulse == nil || *r.Pulse > *ext.highPulse {
			v := *r.Pulse
			ext.highPulse = &v
		}
		if ext.lowPulse == nil || *r.Pulse < *ext.lowPulse {
			v := *r.Pulse
			ext.lowPulse = &v
		}
	}
}

func computeReadingExtremes(sessions []domain.MeasurementSession) readingExtremes {
	ext := readingExtremes{lowSys: 999, lowDia: 999}

	for _, s := range sessions {
		for _, r := range s.Readings {
			ext.trackReading(r)
		}
	}

	return ext
}

func ComputeSummary(sessions []domain.MeasurementSession) Summary {
	if len(sessions) == 0 {
		return Summary{}
	}

	days := make(map[string]bool)
	var allSys, allDia []decimal.Decimal
	var allPulse []decimal.Decimal

	for _, s := range sessions {
		days[s.MeasuredAt.Format("2006-01-02")] = true
		avg := domain.ComputeSessionAverage(s.Readings)
		allSys = append(allSys, avg.AvgSystolic)
		allDia = append(allDia, avg.AvgDiastolic)
		if avg.AvgPulse != nil {
			allPulse = append(allPulse, *avg.AvgPulse)
		}
	}

	ext := computeReadingExtremes(sessions)
	morningAvg, eveningAvg := ComputePeriodAverages(sessions)

	return Summary{
		DaysMeasured: len(days),
		SessionCount: len(sessions),
		ReadingCount: ext.count,
		OverallAvg:   averageFromDecimals(allSys, allDia, allPulse),
		MorningAvg:   morningAvg,
		EveningAvg:   eveningAvg,
		HighestSys:   ext.highSys,
		LowestSys:    ext.lowSys,
		HighestDia:   ext.highDia,
		LowestDia:    ext.lowDia,
		HighestPulse: ext.highPulse,
		LowestPulse:  ext.lowPulse,
	}
}
