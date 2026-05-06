package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"journal/internal/config"
	"journal/internal/journal"
)

// entryInfo holds parsed information about a journal entry
type entryInfo struct {
	existingStr string
	sections    []journal.Section
	rendered    string
	sameTime    bool
	timeValue   string
}

// processEntry handles the main entry processing logic
func processEntry(entryFile string, info entryInfo, settings config.Settings) error {
	if err := ensureDir(entryFile); err != nil {
		return err
	}

	args := flag.Args()
	if len(args) > 0 && !useEditor {
		return appendFromArgs(entryFile, info.sections, info.timeValue, settings.Template, info.sameTime, args)
	}

	if err := prepareEntryFile(entryFile, info); err != nil {
		return err
	}
	return openEditor(entryFile, settings.Editor)
}

// prepareEntryFile prepares the entry file with appropriate content before editing
func prepareEntryFile(entryFile string, info entryInfo) error {
	if info.sameTime {
		return nil
	}

	content := info.existingStr
	if content != "" && !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	if content != "" && !strings.HasSuffix(strings.TrimRight(content, "\n"), "\n\n") {
		content += "\n"
	}
	content += info.rendered + "\n"

	if err := os.WriteFile(entryFile, []byte(content), 0644); err != nil {
		return fmt.Errorf("cannot write entry: %w", err)
	}
	return nil
}

// appendFromArgs appends command-line arguments as content to the entry
func appendFromArgs(entryFile string, sections []journal.Section, timeValue, template string, sameTime bool, args []string) error {
	result := journal.AppendToSections(sections, timeValue, template, sameTime, strings.Join(args, " "))
	if err := os.WriteFile(entryFile, []byte(result), 0644); err != nil {
		return fmt.Errorf("cannot write entry: %w", err)
	}
	return nil
}

// ensureDir ensures the directory for the given file path exists
func ensureDir(filePath string) error {
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return fmt.Errorf("cannot create directory: %w", err)
	}
	return nil
}

// parseEntryFile reads and parses the entry file, returning structured info
func parseEntryFile(entryFile, template, hourFormat string) entryInfo {
	existing, _ := journal.ReadFile(entryFile)
	existingStr := string(existing)
	sections := journal.ParseMarkdown(existingStr)

	now := time.Now()
	var timeValue string
	if hourFormat == "12" {
		timeValue = now.Format("03:04 PM")
	} else {
		timeValue = now.Format("15:04")
	}
	rendered := journal.RenderHeader(template, timeValue)
	sameTime := len(sections) > 0 && sections[len(sections)-1].Heading == rendered

	return entryInfo{
		existingStr: existingStr,
		sections:    sections,
		rendered:    rendered,
		sameTime:    sameTime,
		timeValue:   timeValue,
	}
}
