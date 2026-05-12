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

package domain

import "fmt"

const (
	MinSystolic  = 70
	MaxSystolic  = 250
	MinDiastolic = 40
	MaxDiastolic = 150
	MinPulse     = 30
	MaxPulse     = 220
)

type ValidationError struct {
	Field     string `json:"field"`
	Message   string `json:"message"`
	IsWarning bool   `json:"is_warning"`
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

func ValidateReading(r Reading, index int) []ValidationError {
	var errs []ValidationError

	prefix := fmt.Sprintf("readings[%d]", index)

	if r.Systolic < MinSystolic || r.Systolic > MaxSystolic {
		errs = append(errs, ValidationError{
			Field:   prefix + ".systolic",
			Message: fmt.Sprintf("systolic must be between %d and %d", MinSystolic, MaxSystolic),
		})
	}

	if r.Diastolic < MinDiastolic || r.Diastolic > MaxDiastolic {
		errs = append(errs, ValidationError{
			Field:   prefix + ".diastolic",
			Message: fmt.Sprintf("diastolic must be between %d and %d", MinDiastolic, MaxDiastolic),
		})
	}

	if r.Pulse != nil {
		if *r.Pulse < MinPulse || *r.Pulse > MaxPulse {
			errs = append(errs, ValidationError{
				Field:   prefix + ".pulse",
				Message: fmt.Sprintf("pulse must be between %d and %d", MinPulse, MaxPulse),
			})
		}
	}

	if r.Systolic > 0 && r.Diastolic > 0 && r.Systolic <= r.Diastolic {
		errs = append(errs, ValidationError{
			Field:   prefix + ".systolic",
			Message: "systolic must be greater than diastolic",
		})
	}

	return errs
}

func ValidateSession(s MeasurementSession) []ValidationError {
	var errs []ValidationError

	if s.MeasuredAt.IsZero() {
		errs = append(errs, ValidationError{
			Field:   "measured_at",
			Message: "measurement time is required",
		})
	}

	if s.Period == "" {
		errs = append(errs, ValidationError{
			Field:   "period",
			Message: "period is required",
		})
	} else if s.Period != PeriodMorning && s.Period != PeriodEvening && s.Period != PeriodCustom {
		errs = append(errs, ValidationError{
			Field:   "period",
			Message: "period must be morning, evening, or custom",
		})
	}

	if len(s.Readings) == 0 {
		errs = append(errs, ValidationError{
			Field:   "readings",
			Message: "at least one reading is required",
		})
	}

	for i, r := range s.Readings {
		errs = append(errs, ValidateReading(r, i)...)
	}

	return errs
}
