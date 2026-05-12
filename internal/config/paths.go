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
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

func DataDir() (string, error) {
	if env := os.Getenv("BP_DATA_DIR"); env != "" {
		return env, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolving home directory: %w", err)
	}

	if runtime.GOOS == "darwin" {
		return filepath.Join(home, "Library", "Application Support", "bp-journal"), nil
	}

	return filepath.Join(home, ".local", "share", "bp-journal"), nil
}

func EnsureDataDir(dataDir string) error {
	if err := os.MkdirAll(dataDir, 0o750); err != nil {
		return fmt.Errorf("creating data directory: %w", err)
	}

	if err := os.MkdirAll(SessionsDir(dataDir), 0o750); err != nil {
		return fmt.Errorf("creating sessions directory: %w", err)
	}

	return nil
}

func SessionsDir(dataDir string) string {
	return filepath.Join(dataDir, "sessions")
}

func SettingsPath(dataDir string) string {
	return filepath.Join(dataDir, "settings.toml")
}

func MonthFileName(t time.Time) string {
	return t.Format("2006-01") + ".toml"
}
