package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDataDir(t *testing.T) {
	dir, err := DataDir()
	if err != nil {
		t.Fatalf("DataDir() error: %v", err)
	}
	if dir == "" {
		t.Fatal("DataDir() returned empty string")
	}
}

func TestDataDirEnvOverride(t *testing.T) {
	t.Setenv("BP_DATA_DIR", "/tmp/bp-test")
	dir, err := DataDir()
	if err != nil {
		t.Fatalf("DataDir() error: %v", err)
	}
	if dir != "/tmp/bp-test" {
		t.Errorf("DataDir() = %q, want /tmp/bp-test", dir)
	}
}

func TestEnsureDataDir(t *testing.T) {
	base := t.TempDir()
	dataDir := filepath.Join(base, "bp-journal")

	if err := EnsureDataDir(dataDir); err != nil {
		t.Fatalf("EnsureDataDir() error: %v", err)
	}

	info, err := os.Stat(dataDir)
	if err != nil {
		t.Fatalf("data dir not created: %v", err)
	}
	if !info.IsDir() {
		t.Fatal("data dir is not a directory")
	}

	sessionsDir := filepath.Join(dataDir, "sessions")
	info, err = os.Stat(sessionsDir)
	if err != nil {
		t.Fatalf("sessions dir not created: %v", err)
	}
	if !info.IsDir() {
		t.Fatal("sessions dir is not a directory")
	}
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
		got := MonthFileName(tt.time)
		if got != tt.want {
			t.Errorf("MonthFileName(%v) = %q, want %q", tt.time, got, tt.want)
		}
	}
}
