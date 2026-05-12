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

import "time"

type Period string

const (
	PeriodMorning Period = "morning"
	PeriodEvening Period = "evening"
	PeriodCustom  Period = "custom"
)

type Arm string

const (
	ArmLeft  Arm = "left"
	ArmRight Arm = "right"
)

type Position string

const (
	PositionSitting  Position = "sitting"
	PositionLying    Position = "lying"
	PositionStanding Position = "standing"
)

type Reading struct {
	ID        string    `toml:"id" json:"id"`
	ReadingNo int       `toml:"reading_no" json:"reading_no"`
	Systolic  int       `toml:"systolic" json:"systolic"`
	Diastolic int       `toml:"diastolic" json:"diastolic"`
	Pulse     *int      `toml:"pulse,omitempty" json:"pulse,omitempty"`
	CreatedAt time.Time `toml:"created_at" json:"created_at"`
}

type MeasurementSession struct {
	ID         string    `toml:"id" json:"id"`
	MeasuredAt time.Time `toml:"measured_at" json:"measured_at"`
	Period     Period    `toml:"period" json:"period"`

	Arm          Arm      `toml:"arm,omitempty" json:"arm,omitempty"`
	Position     Position `toml:"position,omitempty" json:"position,omitempty"`
	CuffLocation string   `toml:"cuff_location,omitempty" json:"cuff_location,omitempty"`
	DeviceName   string   `toml:"device_name,omitempty" json:"device_name,omitempty"`

	Symptoms     []string `toml:"symptoms,omitempty" json:"symptoms,omitempty"`
	ContextNotes []string `toml:"context_notes,omitempty" json:"context_notes,omitempty"`
	Notes        string   `toml:"notes,omitempty" json:"notes,omitempty"`

	Readings []Reading `toml:"readings" json:"readings"`

	CreatedAt time.Time `toml:"created_at" json:"created_at"`
	UpdatedAt time.Time `toml:"updated_at" json:"updated_at"`
}

type MonthFile struct {
	Sessions []MeasurementSession `toml:"sessions"`
}

type Settings struct {
	PatientName     string    `toml:"patient_name" json:"patient_name"`
	DefaultArm      Arm       `toml:"default_arm" json:"default_arm"`
	DefaultPosition Position  `toml:"default_position" json:"default_position"`
	DeviceName      string    `toml:"device_name" json:"device_name"`
	UpdatedAt       time.Time `toml:"updated_at" json:"updated_at"`
}

type SettingsFile struct {
	Settings Settings `toml:"settings"`
}
