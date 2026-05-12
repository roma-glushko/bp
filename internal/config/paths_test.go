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
