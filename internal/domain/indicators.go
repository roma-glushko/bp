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
