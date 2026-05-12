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

func ComputeSummary(sessions []domain.MeasurementSession) Summary {
	if len(sessions) == 0 {
		return Summary{}
	}

	days := make(map[string]bool)
	var allSys, allDia []decimal.Decimal
	var allPulse []decimal.Decimal
	var readingCount int

	highSys, lowSys := 0, 999
	highDia, lowDia := 0, 999
	var highPulse, lowPulse *int

	for _, s := range sessions {
		days[s.MeasuredAt.Format("2006-01-02")] = true
		avg := domain.ComputeSessionAverage(s.Readings)
		allSys = append(allSys, avg.AvgSystolic)
		allDia = append(allDia, avg.AvgDiastolic)
		if avg.AvgPulse != nil {
			allPulse = append(allPulse, *avg.AvgPulse)
		}

		for _, r := range s.Readings {
			readingCount++

			if r.Systolic > highSys {
				highSys = r.Systolic
			}

			if r.Systolic < lowSys {
				lowSys = r.Systolic
			}

			if r.Diastolic > highDia {
				highDia = r.Diastolic
			}

			if r.Diastolic < lowDia {
				lowDia = r.Diastolic
			}

			if r.Pulse != nil {
				if highPulse == nil || *r.Pulse > *highPulse {
					v := *r.Pulse
					highPulse = &v
				}
				if lowPulse == nil || *r.Pulse < *lowPulse {
					v := *r.Pulse
					lowPulse = &v
				}
			}
		}
	}

	morningAvg, eveningAvg := ComputePeriodAverages(sessions)

	return Summary{
		DaysMeasured: len(days),
		SessionCount: len(sessions),
		ReadingCount: readingCount,
		OverallAvg:   averageFromDecimals(allSys, allDia, allPulse),
		MorningAvg:   morningAvg,
		EveningAvg:   eveningAvg,
		HighestSys:   highSys,
		LowestSys:    lowSys,
		HighestDia:   highDia,
		LowestDia:    lowDia,
		HighestPulse: highPulse,
		LowestPulse:  lowPulse,
	}
}
