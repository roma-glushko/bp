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
	"fmt"
	"sort"
	"time"

	"github.com/shopspring/decimal"

	"github.com/roma-glushko/bp/internal/domain"
)

type PeriodAverage struct {
	AvgSystolic  decimal.Decimal    `json:"avg_systolic"`
	AvgDiastolic decimal.Decimal    `json:"avg_diastolic"`
	AvgPulse     *decimal.Decimal   `json:"avg_pulse,omitempty"`
	Indicator    domain.BPIndicator `json:"indicator"`
	Count        int                `json:"count"`
}

type DayAverage struct {
	Date     string        `json:"date"`
	Average  PeriodAverage `json:"average"`
	Sessions int           `json:"sessions"`
	Notes    []string      `json:"notes,omitempty"`
}

type WeekAverage struct {
	Week    string        `json:"week"`
	Average PeriodAverage `json:"average"`
	Days    int           `json:"days"`
}

type MonthAverage struct {
	Month   string        `json:"month"`
	Average PeriodAverage `json:"average"`
	Days    int           `json:"days"`
}

type SessionEntry struct {
	domain.MeasurementSession
	Average domain.SessionAverage `json:"average"`
}

func averageFromDecimals(sysList, diaList []decimal.Decimal, pulseList []decimal.Decimal) PeriodAverage {
	if len(sysList) == 0 {
		return PeriodAverage{}
	}

	n := decimal.NewFromInt(int64(len(sysList)))

	sumSys := decimal.Zero
	sumDia := decimal.Zero

	for i := range sysList {
		sumSys = sumSys.Add(sysList[i])
		sumDia = sumDia.Add(diaList[i])
	}

	avgSys := sumSys.Div(n)
	avgDia := sumDia.Div(n)

	avg := PeriodAverage{
		AvgSystolic:  avgSys,
		AvgDiastolic: avgDia,
		Indicator:    domain.ClassifyBP(int(avgSys.Round(0).IntPart()), int(avgDia.Round(0).IntPart())),
		Count:        len(sysList),
	}

	if len(pulseList) > 0 {
		sumPulse := decimal.Zero
		for _, p := range pulseList {
			sumPulse = sumPulse.Add(p)
		}
		p := sumPulse.Div(decimal.NewFromInt(int64(len(pulseList))))
		avg.AvgPulse = &p
	}

	return avg
}

func ComputeDayAverages(sessions []domain.MeasurementSession) []DayAverage {
	type dayData struct {
		sysList    []decimal.Decimal
		diaList    []decimal.Decimal
		pulseList  []decimal.Decimal
		sessionCnt int
		notes      []string
	}

	days := make(map[string]*dayData)

	var dayOrder []string

	for _, s := range sessions {
		key := s.MeasuredAt.Format("2006-01-02")
		dd, ok := days[key]
		if !ok {
			dd = &dayData{}
			days[key] = dd
			dayOrder = append(dayOrder, key)
		}

		avg := domain.ComputeSessionAverage(s.Readings)
		dd.sysList = append(dd.sysList, avg.AvgSystolic)
		dd.diaList = append(dd.diaList, avg.AvgDiastolic)
		if avg.AvgPulse != nil {
			dd.pulseList = append(dd.pulseList, *avg.AvgPulse)
		}
		dd.sessionCnt++
		if s.Notes != "" {
			dd.notes = append(dd.notes, s.Notes)
		}
	}

	sort.Strings(dayOrder)

	result := make([]DayAverage, 0, len(dayOrder))
	for _, key := range dayOrder {
		dd := days[key]
		result = append(result, DayAverage{
			Date:     key,
			Average:  averageFromDecimals(dd.sysList, dd.diaList, dd.pulseList),
			Sessions: dd.sessionCnt,
			Notes:    dd.notes,
		})
	}

	return result
}

func ComputeWeekAverages(dayAvgs []DayAverage) []WeekAverage {
	type weekData struct {
		sysList   []decimal.Decimal
		diaList   []decimal.Decimal
		pulseList []decimal.Decimal
		days      int
	}

	weeks := make(map[string]*weekData)
	var weekOrder []string

	for _, da := range dayAvgs {
		t, _ := time.Parse("2006-01-02", da.Date)
		year, week := t.ISOWeek()
		key := fmt.Sprintf("%d-W%02d", year, week)

		wd, ok := weeks[key]
		if !ok {
			wd = &weekData{}

			weeks[key] = wd

			weekOrder = append(weekOrder, key)
		}

		wd.sysList = append(wd.sysList, da.Average.AvgSystolic)
		wd.diaList = append(wd.diaList, da.Average.AvgDiastolic)
		if da.Average.AvgPulse != nil {
			wd.pulseList = append(wd.pulseList, *da.Average.AvgPulse)
		}

		wd.days++
	}

	sort.Strings(weekOrder)

	result := make([]WeekAverage, 0, len(weekOrder))

	for _, key := range weekOrder {
		wd := weeks[key]
		result = append(result, WeekAverage{
			Week:    key,
			Average: averageFromDecimals(wd.sysList, wd.diaList, wd.pulseList),
			Days:    wd.days,
		})
	}

	return result
}

func ComputeMonthAverages(dayAvgs []DayAverage) []MonthAverage {
	type monthData struct {
		sysList   []decimal.Decimal
		diaList   []decimal.Decimal
		pulseList []decimal.Decimal
		days      int
	}

	months := make(map[string]*monthData)

	var monthOrder []string

	for _, da := range dayAvgs {
		key := da.Date[:7]
		md, ok := months[key]

		if !ok {
			md = &monthData{}
			months[key] = md

			monthOrder = append(monthOrder, key)
		}

		md.sysList = append(md.sysList, da.Average.AvgSystolic)
		md.diaList = append(md.diaList, da.Average.AvgDiastolic)

		if da.Average.AvgPulse != nil {
			md.pulseList = append(md.pulseList, *da.Average.AvgPulse)
		}

		md.days++
	}

	sort.Strings(monthOrder)

	result := make([]MonthAverage, 0, len(monthOrder))

	for _, key := range monthOrder {
		md := months[key]
		result = append(result, MonthAverage{
			Month:   key,
			Average: averageFromDecimals(md.sysList, md.diaList, md.pulseList),
			Days:    md.days,
		})
	}

	return result
}

func ComputePeriodAverages(sessions []domain.MeasurementSession) (morning, evening PeriodAverage) {
	var mSys, mDia, mPulse []decimal.Decimal

	var eSys, eDia, ePulse []decimal.Decimal

	for _, s := range sessions {
		avg := domain.ComputeSessionAverage(s.Readings)
		switch s.Period {
		case domain.PeriodMorning:
			mSys = append(mSys, avg.AvgSystolic)
			mDia = append(mDia, avg.AvgDiastolic)
			if avg.AvgPulse != nil {
				mPulse = append(mPulse, *avg.AvgPulse)
			}
		case domain.PeriodEvening:
			eSys = append(eSys, avg.AvgSystolic)
			eDia = append(eDia, avg.AvgDiastolic)
			if avg.AvgPulse != nil {
				ePulse = append(ePulse, *avg.AvgPulse)
			}
		case domain.PeriodCustom:
			// custom-period sessions are excluded from morning/evening split
		}
	}

	return averageFromDecimals(mSys, mDia, mPulse), averageFromDecimals(eSys, eDia, ePulse)
}
