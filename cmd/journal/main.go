package main

import (
	"fmt"
	"os"

	"journal/internal/journal"
)

func main() {
	settings, basePath := loadSettings()
	if basePath == "" {
		return // showRoot was requested
	}

	entryFile, _ := journal.EntryFile(basePath, settings.HourFormat, settings.FileNameFormat)
	info := parseEntryFile(entryFile, settings.Template, settings.HourFormat)

	if err := processEntry(entryFile, info, settings); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Entry saved to %s\n", entryFile)
	saveSettingsIfNeeded(settings)
}
