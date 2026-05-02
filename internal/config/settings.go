package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// Settings は設定ファイルの構造
type Settings struct {
	HourFormat         string `json:"hour_format"`
	Path               string `json:"path"`
	Editor             string `json:"editor"`
	Template           string `json:"template"`
	DefaultProject     bool   `json:"default_project"`
	FileNameFormat     string `json:"file_name_format"`
	AutoCreateSettings bool   `json:"auto_create_settings"`
}

// ConfigPath は設定ファイルのパスを返す
func ConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".journal", "settings.json"), nil
}

func DefaultSettings() Settings {
	home, _ := os.UserHomeDir()
	return Settings{
		HourFormat: "24", Path: filepath.Join(home, ".journal"),
		Editor: "", Template: "## {time}\n",
		DefaultProject: false, FileNameFormat: "02.md",
		AutoCreateSettings: true,
	}
}

func LoadSettings() (Settings, error) {
	s := DefaultSettings()
	configPath, err := ConfigPath()
	if err != nil {
		return s, err
	}
	data, err := os.ReadFile(configPath)
	if err != nil {
		return s, nil
	}
	if err := json.Unmarshal(data, &s); err != nil {
		return DefaultSettings(), nil
	}
	if s.Path == "" {
		home, _ := os.UserHomeDir()
		s.Path = filepath.Join(home, ".journal")
	}
	if strings.HasPrefix(s.Path, "~/") {
		home, _ := os.UserHomeDir()
		s.Path = filepath.Join(home, strings.TrimPrefix(s.Path, "~/"))
	}
	return s, nil
}

func SaveSettings(s Settings) error {
	configPath, err := ConfigPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, data, 0644)
}

func Exists() bool {
	p, err := ConfigPath()
	if err != nil {
		return false
	}
	_, err = os.Stat(p)
	return err == nil
}
