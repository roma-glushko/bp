package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func intPtr(v int) *int { return &v }

func TestValidateReading(t *testing.T) {
	tests := []struct {
		name     string
		reading  Reading
		wantErrs int
	}{
		{
			name:     "valid reading with pulse",
			reading:  Reading{Systolic: 120, Diastolic: 80, Pulse: intPtr(75)},
			wantErrs: 0,
		},
		{
			name:     "valid reading without pulse",
			reading:  Reading{Systolic: 120, Diastolic: 80, Pulse: nil},
			wantErrs: 0,
		},
		{
			name:     "systolic too low",
			reading:  Reading{Systolic: 50, Diastolic: 80},
			wantErrs: 2, // out of range + less than diastolic
		},
		{
			name:     "systolic too high",
			reading:  Reading{Systolic: 260, Diastolic: 80},
			wantErrs: 1,
		},
		{
			name:     "diastolic too low",
			reading:  Reading{Systolic: 120, Diastolic: 30},
			wantErrs: 1,
		},
		{
			name:     "diastolic too high",
			reading:  Reading{Systolic: 120, Diastolic: 160},
			wantErrs: 2, // out of range + systolic less than diastolic
		},
		{
			name:     "pulse too low",
			reading:  Reading{Systolic: 120, Diastolic: 80, Pulse: intPtr(20)},
			wantErrs: 1,
		},
		{
			name:     "pulse too high",
			reading:  Reading{Systolic: 120, Diastolic: 80, Pulse: intPtr(230)},
			wantErrs: 1,
		},
		{
			name:     "systolic less than diastolic",
			reading:  Reading{Systolic: 80, Diastolic: 90},
			wantErrs: 1,
		},
		{
			name:     "boundary valid low",
			reading:  Reading{Systolic: MinSystolic, Diastolic: MinDiastolic, Pulse: intPtr(MinPulse)},
			wantErrs: 0,
		},
		{
			name:     "boundary valid high",
			reading:  Reading{Systolic: MaxSystolic, Diastolic: MaxDiastolic, Pulse: intPtr(MaxPulse)},
			wantErrs: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := ValidateReading(tt.reading, 0)
			require.Len(t, errs, tt.wantErrs)
		})
	}
}

func TestValidateSession(t *testing.T) {
	validReading := Reading{Systolic: 120, Diastolic: 80, Pulse: intPtr(75)}

	tests := []struct {
		name     string
		session  MeasurementSession
		wantErrs int
	}{
		{
			name: "valid session",
			session: MeasurementSession{
				MeasuredAt: time.Now(),
				Period:     PeriodMorning,
				Readings:   []Reading{validReading},
			},
			wantErrs: 0,
		},
		{
			name: "missing measured_at",
			session: MeasurementSession{
				Period:   PeriodMorning,
				Readings: []Reading{validReading},
			},
			wantErrs: 1,
		},
		{
			name: "missing period",
			session: MeasurementSession{
				MeasuredAt: time.Now(),
				Readings:   []Reading{validReading},
			},
			wantErrs: 1,
		},
		{
			name: "invalid period",
			session: MeasurementSession{
				MeasuredAt: time.Now(),
				Period:     "afternoon",
				Readings:   []Reading{validReading},
			},
			wantErrs: 1,
		},
		{
			name: "no readings",
			session: MeasurementSession{
				MeasuredAt: time.Now(),
				Period:     PeriodMorning,
				Readings:   []Reading{},
			},
			wantErrs: 1,
		},
		{
			name: "invalid reading in session",
			session: MeasurementSession{
				MeasuredAt: time.Now(),
				Period:     PeriodMorning,
				Readings:   []Reading{{Systolic: 50, Diastolic: 80}},
			},
			wantErrs: 2, // out of range + systolic less than diastolic
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := ValidateSession(tt.session)
			require.Len(t, errs, tt.wantErrs)
		})
	}
}
