package domain

type BPStatus string

const (
	BPStatusGreat    BPStatus = "great"
	BPStatusElevated BPStatus = "elevated"
	BPStatusHigh     BPStatus = "high"
	BPStatusVeryHigh BPStatus = "very_high"
)

type BPIndicator struct {
	Status BPStatus `json:"status"`
	Color  string   `json:"color"`
	Label  string   `json:"label"`
}

func ClassifyBP(systolic, diastolic int) BPIndicator {
	if systolic >= 180 || diastolic >= 120 {
		return BPIndicator{Status: BPStatusVeryHigh, Color: "darkred", Label: "Very high"}
	}
	if systolic >= 140 || diastolic >= 90 {
		return BPIndicator{Status: BPStatusHigh, Color: "red", Label: "Too high"}
	}
	if systolic >= 120 || diastolic >= 80 {
		return BPIndicator{Status: BPStatusElevated, Color: "yellow", Label: "Elevated"}
	}
	return BPIndicator{Status: BPStatusGreat, Color: "green", Label: "Great"}
}
