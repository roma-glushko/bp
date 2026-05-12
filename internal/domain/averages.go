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

import "github.com/shopspring/decimal"

type SessionAverage struct {
	AvgSystolic  decimal.Decimal  `json:"avg_systolic"`
	AvgDiastolic decimal.Decimal  `json:"avg_diastolic"`
	AvgPulse     *decimal.Decimal `json:"avg_pulse,omitempty"`
	Indicator    BPIndicator      `json:"indicator"`
}

func ComputeSessionAverage(readings []Reading) SessionAverage {
	if len(readings) == 0 {
		return SessionAverage{}
	}

	sumSys := decimal.Zero
	sumDia := decimal.Zero
	sumPulse := decimal.Zero
	pulseCount := 0

	for _, r := range readings {
		sumSys = sumSys.Add(decimal.NewFromInt(int64(r.Systolic)))
		sumDia = sumDia.Add(decimal.NewFromInt(int64(r.Diastolic)))
		if r.Pulse != nil {
			sumPulse = sumPulse.Add(decimal.NewFromInt(int64(*r.Pulse)))
			pulseCount++
		}
	}

	n := decimal.NewFromInt(int64(len(readings)))
	avgSys := sumSys.Div(n)
	avgDia := sumDia.Div(n)

	roundedSys := avgSys.Round(0).IntPart()
	roundedDia := avgDia.Round(0).IntPart()

	avg := SessionAverage{
		AvgSystolic:  avgSys,
		AvgDiastolic: avgDia,
		Indicator:    ClassifyBP(int(roundedSys), int(roundedDia)),
	}

	if pulseCount > 0 {
		p := sumPulse.Div(decimal.NewFromInt(int64(pulseCount)))
		avg.AvgPulse = &p
	}

	return avg
}
