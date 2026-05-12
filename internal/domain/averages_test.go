package domain

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestComputeSessionAverage(t *testing.T) {
	pulse78 := 78
	pulse81 := 81

	readings := []Reading{
		{Systolic: 126, Diastolic: 86, Pulse: &pulse78},
		{Systolic: 119, Diastolic: 83, Pulse: &pulse81},
	}

	avg := ComputeSessionAverage(readings)

	assert.True(t, avg.AvgSystolic.Equal(decimal.NewFromFloat(122.5)), "AvgSystolic = %s, want 122.5", avg.AvgSystolic)
	assert.True(t, avg.AvgDiastolic.Equal(decimal.NewFromFloat(84.5)), "AvgDiastolic = %s, want 84.5", avg.AvgDiastolic)

	require.NotNil(t, avg.AvgPulse)
	assert.True(t, avg.AvgPulse.Equal(decimal.NewFromFloat(79.5)), "AvgPulse = %s, want 79.5", avg.AvgPulse)

	assert.Equal(t, BPStatusElevated, avg.Indicator.Status)
}

func TestComputeSessionAverageThreeReadings(t *testing.T) {
	pulse := 70

	readings := []Reading{
		{Systolic: 120, Diastolic: 80, Pulse: &pulse},
		{Systolic: 121, Diastolic: 81, Pulse: &pulse},
		{Systolic: 122, Diastolic: 82, Pulse: &pulse},
	}

	avg := ComputeSessionAverage(readings)

	assert.True(t, avg.AvgSystolic.Equal(decimal.NewFromInt(121)), "AvgSystolic = %s, want 121", avg.AvgSystolic)
	assert.True(t, avg.AvgDiastolic.Equal(decimal.NewFromInt(81)), "AvgDiastolic = %s, want 81", avg.AvgDiastolic)
}

func TestComputeSessionAverageNoPulse(t *testing.T) {
	readings := []Reading{
		{Systolic: 115, Diastolic: 75},
		{Systolic: 110, Diastolic: 70},
	}

	avg := ComputeSessionAverage(readings)

	assert.Nil(t, avg.AvgPulse)
	assert.Equal(t, BPStatusGreat, avg.Indicator.Status)
}

func TestComputeSessionAverageEmpty(t *testing.T) {
	avg := ComputeSessionAverage(nil)

	assert.True(t, avg.AvgSystolic.IsZero())
	assert.True(t, avg.AvgDiastolic.IsZero())
}
