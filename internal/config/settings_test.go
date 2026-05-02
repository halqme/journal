package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultSettings(t *testing.T) {
	s := DefaultSettings()
	if s.HourFormat != "24" {
		t.Errorf("expected hour_format '24', got %q", s.HourFormat)
	}
	if s.Template != "## {time}\n" {
		t.Errorf("expected default template, got %q", s.Template)
	}
	if !s.AutoCreateSettings {
		t.Error("expected AutoCreateSettings to be true")
	}
}

func TestSettings_JSONRoundTrip(t *testing.T) {
	s := Settings{
		HourFormat:         "12",
		Path:               "/custom/path",
		Editor:             "nvim",
		Template:           "## {time}\n---\n",
		DefaultProject:     true,
		FileNameFormat:     "2006-01-02.md",
		AutoCreateSettings: false,
	}

	data, err := json.Marshal(s)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var decoded Settings
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if decoded.HourFormat != s.HourFormat {
		t.Errorf("hour_format mismatch: got %q", decoded.HourFormat)
	}
	if decoded.Path != s.Path {
		t.Errorf("path mismatch: got %q", decoded.Path)
	}
	if decoded.Editor != s.Editor {
		t.Errorf("editor mismatch: got %q", decoded.Editor)
	}
	if decoded.DefaultProject != s.DefaultProject {
		t.Errorf("default_project mismatch: got %v", decoded.DefaultProject)
	}
}

func TestLoadSettings_MissingFile(t *testing.T) {
	s, err := LoadSettings()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.HourFormat != "24" {
		t.Errorf("expected default hour_format, got %q", s.HourFormat)
	}
}

func TestLoadSettings_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".journal", "settings.json")
	os.MkdirAll(filepath.Dir(configPath), 0755)
	os.WriteFile(configPath, []byte("not valid json{{{"), 0644)

	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", oldHome)

	s, err := LoadSettings()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.HourFormat != "24" {
		t.Errorf("expected default hour_format, got %q", s.HourFormat)
	}
}

func TestSaveSettings(t *testing.T) {
	tmpDir := t.TempDir()
	os.Setenv("HOME", tmpDir)

	s := Settings{
		HourFormat: "12",
		Path:       filepath.Join(tmpDir, ".journal"),
		Template:   "## {time}\n",
	}

	if err := SaveSettings(s); err != nil {
		t.Fatalf("failed to save settings: %v", err)
	}

	configPath := filepath.Join(tmpDir, ".journal", "settings.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read saved file: %v", err)
	}

	var loaded Settings
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Fatalf("failed to unmarshal saved file: %v", err)
	}
	if loaded.HourFormat != "12" {
		t.Errorf("expected hour_format '12', got %q", loaded.HourFormat)
	}
}
