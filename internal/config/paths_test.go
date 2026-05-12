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

package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDataDir(t *testing.T) {
	dir, err := DataDir()

	require.NoError(t, err)
	require.NotEmpty(t, dir)
}

func TestDataDirEnvOverride(t *testing.T) {
	t.Setenv("BP_DATA_DIR", "/tmp/bp-test")

	dir, err := DataDir()

	require.NoError(t, err)
	require.Equal(t, "/tmp/bp-test", dir)
}

func TestEnsureDataDir(t *testing.T) {
	base := t.TempDir()
	dataDir := filepath.Join(base, "bp-journal")

	require.NoError(t, EnsureDataDir(dataDir))

	info, err := os.Stat(dataDir)

	require.NoError(t, err)
	require.True(t, info.IsDir())

	sessionsDir := filepath.Join(dataDir, "sessions")
	info, err = os.Stat(sessionsDir)

	require.NoError(t, err)
	require.True(t, info.IsDir())
}

func TestMonthFileName(t *testing.T) {
	tests := []struct {
		time time.Time
		want string
	}{
		{time.Date(2026, 5, 7, 9, 0, 0, 0, time.UTC), "2026-05.toml"},
		{time.Date(2026, 12, 31, 23, 59, 0, 0, time.UTC), "2026-12.toml"},
		{time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC), "2027-01.toml"},
	}

	for _, tt := range tests {
		require.Equal(t, tt.want, MonthFileName(tt.time))
	}
}
