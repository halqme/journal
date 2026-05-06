package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"journal/internal/config"
)

var (
	useEditor  bool
	useProject bool
	showRoot   bool
)

func init() {
	flag.BoolVar(&useEditor, "e", false, "open $EDITOR")
	flag.BoolVar(&useEditor, "editor", false, "open $EDITOR")
	flag.BoolVar(&useProject, "p", false, "create entry in current directory (.journal/entries)")
	flag.BoolVar(&useProject, "project", false, "create entry in current directory (.journal/entries)")
	flag.BoolVar(&showRoot, "root", false, "print journal base path")
}

// loadSettings loads configuration and returns settings + basePath (empty string if showRoot)
func loadSettings() (config.Settings, string) {
	flag.Parse()

	settings, err := config.LoadSettings()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	basePath := resolveBasePath(settings, useProject)

	if showRoot {
		fmt.Println(basePath)
		return settings, ""
	}

	return settings, basePath
}

// resolveBasePath determines the base path for journal entries
func resolveBasePath(settings config.Settings, useProject bool) string {
	basePath := settings.Path
	if useProject || settings.DefaultProject {
		cwd, err := os.Getwd()
		if err != nil {
			return basePath // fallback to default
		}
		basePath = filepath.Join(cwd, ".journal", "entries")
	}
	return basePath
}

// saveSettingsIfNeeded saves settings if auto-create is enabled
func saveSettingsIfNeeded(settings config.Settings) {
	if settings.AutoCreateSettings && !config.Exists() {
		if err := config.SaveSettings(settings); err != nil {
			fmt.Fprintf(os.Stderr, "warning: could not save settings: %v\n", err)
		}
	}
}
