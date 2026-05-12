package domain

import "github.com/shopspring/decimal"

type SessionAverage struct {
	AvgSystolic  decimal.Decimal `json:"avg_systolic"`
	AvgDiastolic decimal.Decimal `json:"avg_diastolic"`
	AvgPulse     *decimal.Decimal `json:"avg_pulse,omitempty"`
	Indicator    BPIndicator     `json:"indicator"`
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
