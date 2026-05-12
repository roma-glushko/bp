package report

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/roma-glushko/bp/internal/domain"
)

func intPtr(v int) *int { return &v }

func TestComputeDayAverages(t *testing.T) {
	sessions := []domain.MeasurementSession{
		{
			MeasuredAt: time.Date(2026, 5, 7, 9, 0, 0, 0, time.UTC),
			Period:     domain.PeriodMorning,
			Readings: []domain.Reading{
				{Systolic: 126, Diastolic: 86, Pulse: intPtr(78)},
				{Systolic: 119, Diastolic: 83, Pulse: intPtr(81)},
			},
		},
		{
			MeasuredAt: time.Date(2026, 5, 7, 22, 0, 0, 0, time.UTC),
			Period:     domain.PeriodEvening,
			Readings: []domain.Reading{
				{Systolic: 119, Diastolic: 78, Pulse: intPtr(75)},
				{Systolic: 119, Diastolic: 80, Pulse: intPtr(78)},
			},
		},
	}

	dayAvgs := ComputeDayAverages(sessions)
	require.Len(t, dayAvgs, 1)

	day := dayAvgs[0]
	assert.Equal(t, "2026-05-07", day.Date)
	assert.Equal(t, 2, day.Sessions)

	assert.True(t, day.Average.AvgSystolic.Round(0).Equal(decimal.NewFromInt(121)))
	assert.True(t, day.Average.AvgDiastolic.Round(0).Equal(decimal.NewFromInt(82)))
	require.NotNil(t, day.Average.AvgPulse)
	assert.True(t, day.Average.AvgPulse.Round(0).Equal(decimal.NewFromInt(78)))
}

func TestComputeWeekAverages(t *testing.T) {
	dayAvgs := []DayAverage{
		{
			Date:    "2026-05-04",
			Average: PeriodAverage{AvgSystolic: decimal.NewFromInt(120), AvgDiastolic: decimal.NewFromInt(80)},
		},
		{
			Date:    "2026-05-05",
			Average: PeriodAverage{AvgSystolic: decimal.NewFromInt(130), AvgDiastolic: decimal.NewFromInt(85)},
		},
		{
			Date:    "2026-05-12",
			Average: PeriodAverage{AvgSystolic: decimal.NewFromInt(115), AvgDiastolic: decimal.NewFromInt(75)},
		},
	}

	weekAvgs := ComputeWeekAverages(dayAvgs)
	require.Len(t, weekAvgs, 2)

	assert.Equal(t, "2026-W19", weekAvgs[0].Week)
	assert.Equal(t, 2, weekAvgs[0].Days)
	assert.True(t, weekAvgs[0].Average.AvgSystolic.Equal(decimal.NewFromInt(125)))

	assert.Equal(t, "2026-W20", weekAvgs[1].Week)
	assert.Equal(t, 1, weekAvgs[1].Days)
}

func TestComputeMonthAverages(t *testing.T) {
	dayAvgs := []DayAverage{
		{
			Date:    "2026-05-07",
			Average: PeriodAverage{AvgSystolic: decimal.NewFromInt(121), AvgDiastolic: decimal.NewFromInt(82)},
		},
		{
			Date:    "2026-05-08",
			Average: PeriodAverage{AvgSystolic: decimal.NewFromInt(118), AvgDiastolic: decimal.NewFromInt(78)},
		},
		{
			Date:    "2026-06-01",
			Average: PeriodAverage{AvgSystolic: decimal.NewFromInt(115), AvgDiastolic: decimal.NewFromInt(75)},
		},
	}

	monthAvgs := ComputeMonthAverages(dayAvgs)
	require.Len(t, monthAvgs, 2)

	assert.Equal(t, "2026-05", monthAvgs[0].Month)
	assert.Equal(t, 2, monthAvgs[0].Days)
	assert.True(t, monthAvgs[0].Average.AvgSystolic.Round(0).Equal(decimal.NewFromInt(120)))

	assert.Equal(t, "2026-06", monthAvgs[1].Month)
	assert.Equal(t, 1, monthAvgs[1].Days)
}

func TestComputePeriodAverages(t *testing.T) {
	sessions := []domain.MeasurementSession{
		{
			Period:   domain.PeriodMorning,
			Readings: []domain.Reading{{Systolic: 126, Diastolic: 86, Pulse: intPtr(78)}},
		},
		{
			Period:   domain.PeriodMorning,
			Readings: []domain.Reading{{Systolic: 120, Diastolic: 80, Pulse: intPtr(72)}},
		},
		{
			Period:   domain.PeriodEvening,
			Readings: []domain.Reading{{Systolic: 119, Diastolic: 78, Pulse: intPtr(75)}},
		},
	}

	morning, evening := ComputePeriodAverages(sessions)

	assert.Equal(t, 2, morning.Count)
	assert.True(t, morning.AvgSystolic.Equal(decimal.NewFromInt(123)))
	assert.True(t, morning.AvgDiastolic.Equal(decimal.NewFromInt(83)))

	assert.Equal(t, 1, evening.Count)
	assert.True(t, evening.AvgSystolic.Equal(decimal.NewFromInt(119)))
}

func TestComputeSummary(t *testing.T) {
	sessions := []domain.MeasurementSession{
		{
			MeasuredAt: time.Date(2026, 5, 7, 9, 0, 0, 0, time.UTC),
			Period:     domain.PeriodMorning,
			Readings: []domain.Reading{
				{Systolic: 126, Diastolic: 86, Pulse: intPtr(78)},
				{Systolic: 119, Diastolic: 83, Pulse: intPtr(81)},
			},
		},
		{
			MeasuredAt: time.Date(2026, 5, 7, 22, 0, 0, 0, time.UTC),
			Period:     domain.PeriodEvening,
			Readings: []domain.Reading{
				{Systolic: 119, Diastolic: 78, Pulse: intPtr(75)},
			},
		},
		{
			MeasuredAt: time.Date(2026, 5, 8, 9, 0, 0, 0, time.UTC),
			Period:     domain.PeriodMorning,
			Readings: []domain.Reading{
				{Systolic: 115, Diastolic: 75, Pulse: intPtr(70)},
			},
		},
	}

	s := ComputeSummary(sessions)

	assert.Equal(t, 2, s.DaysMeasured)
	assert.Equal(t, 3, s.SessionCount)
	assert.Equal(t, 4, s.ReadingCount)
	assert.Equal(t, 126, s.HighestSys)
	assert.Equal(t, 115, s.LowestSys)
	assert.Equal(t, 86, s.HighestDia)
	assert.Equal(t, 75, s.LowestDia)
	require.NotNil(t, s.HighestPulse)
	assert.Equal(t, 81, *s.HighestPulse)
	require.NotNil(t, s.LowestPulse)
	assert.Equal(t, 70, *s.LowestPulse)
}

func TestComputeSummaryEmpty(t *testing.T) {
	s := ComputeSummary(nil)
	assert.Equal(t, 0, s.SessionCount)
}
