package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"journal/internal/config"
	"journal/internal/journal"
)

func main() {
	var useEditor, useProject, showRoot bool
	flag.BoolVar(&useEditor, "e", false, "open $EDITOR")
	flag.BoolVar(&useEditor, "editor", false, "open $EDITOR")
	flag.BoolVar(&useProject, "p", false, "create entry in current directory (.journal/entries)")
	flag.BoolVar(&useProject, "project", false, "create entry in current directory (.journal/entries)")
	flag.BoolVar(&showRoot, "root", false, "print journal base path")
	flag.Parse()

	settings, err := config.LoadSettings()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	basePath := settings.Path
	if useProject || settings.DefaultProject {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: cannot get current directory: %v\n", err)
			os.Exit(1)
		}
		basePath = filepath.Join(cwd, ".journal", "entries")
	}

	if showRoot {
		fmt.Println(basePath)
		return
	}

	entryFile, timeValue := journal.EntryFile(basePath, settings.HourFormat, settings.FileNameFormat)

	if err := os.MkdirAll(filepath.Dir(entryFile), 0755); err != nil {
		fmt.Fprintf(os.Stderr, "error: cannot create directory: %v\n", err)
		os.Exit(1)
	}

	existing, _ := journal.ReadFile(entryFile)
	existingStr := string(existing)
	sections := journal.ParseMarkdown(existingStr)

	rendered := journal.RenderHeader(settings.Template, timeValue)
	sameTime := len(sections) > 0 && sections[len(sections)-1].Heading == rendered

	var result string
	args := flag.Args()

	if len(args) > 0 && !useEditor {
		result = journal.AppendToSections(sections, timeValue, settings.Template, sameTime, strings.Join(args, " "))
	} else {
		var tmpContent string
		if sameTime {
			tmpContent = existingStr
		} else {
			tmpContent = existingStr + rendered + "\n"
		}
		result, err = journal.OpenEditor(tmpContent, settings.Editor)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	}

	if err := os.WriteFile(entryFile, []byte(result), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "error: cannot write entry: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Entry saved to %s\n", entryFile)

	if settings.AutoCreateSettings && !config.Exists() {
		if err := config.SaveSettings(settings); err != nil {
			fmt.Fprintf(os.Stderr, "warning: could not save settings: %v\n", err)
		}
	}
}
