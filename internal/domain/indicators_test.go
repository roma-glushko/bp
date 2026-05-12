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

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClassifyBP(t *testing.T) {
	tests := []struct {
		name      string
		systolic  int
		diastolic int
		want      BPStatus
	}{
		{"119/78 is great", 119, 78, BPStatusGreat},
		{"110/70 is great", 110, 70, BPStatusGreat},
		{"126/76 elevated by systolic", 126, 76, BPStatusElevated},
		{"118/86 elevated by diastolic", 118, 86, BPStatusElevated},
		{"139/89 elevated upper bound", 139, 89, BPStatusElevated},
		{"142/82 high by systolic", 142, 82, BPStatusHigh},
		{"135/92 high by diastolic", 135, 92, BPStatusHigh},
		{"160/100 high both", 160, 100, BPStatusHigh},
		{"181/95 very high by systolic", 181, 95, BPStatusVeryHigh},
		{"150/121 very high by diastolic", 150, 121, BPStatusVeryHigh},
		{"200/130 very high both", 200, 130, BPStatusVeryHigh},
		// boundary cases
		{"120/79 elevated at systolic boundary", 120, 79, BPStatusElevated},
		{"119/80 elevated at diastolic boundary", 119, 80, BPStatusElevated},
		{"140/89 high at systolic boundary", 140, 89, BPStatusHigh},
		{"139/90 high at diastolic boundary", 139, 90, BPStatusHigh},
		{"180/89 very high at systolic boundary", 180, 89, BPStatusVeryHigh},
		{"179/120 very high at diastolic boundary", 179, 120, BPStatusVeryHigh},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ClassifyBP(tt.systolic, tt.diastolic)
			assert.Equal(t, tt.want, got.Status)
		})
	}
}
