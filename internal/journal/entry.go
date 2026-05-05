package journal

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func EntryFile(basePath, hourFormat, fileNameFormat string) (file, timeValue string) {
	now := time.Now()
	var timeStr string
	if hourFormat == "12" {
		timeStr = now.Format("03:04 PM")
	} else {
		timeStr = now.Format("15:04")
	}
	entryDir := filepath.Join(basePath, now.Format("2006"), now.Format("01"))
	return filepath.Join(entryDir, now.Format(fileNameFormat)), timeStr
}

func ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func CreateTempFile(content string) (string, error) {
	f, err := os.CreateTemp("", "journal-*.md")
	if err != nil {
		return "", fmt.Errorf("create temp file: %w", err)
	}
	path := f.Name()
	if _, err := f.WriteString(content); err != nil {
		f.Close()
		os.Remove(path)
		return "", fmt.Errorf("write temp file: %w", err)
	}
	f.Close()
	return path, nil
}

func OpenEditor(content, editorOverride string) (string, error) {
	tmpPath, err := CreateTempFile(content)
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpPath)

	editor := editorOverride
	if editor == "" {
		editor = os.Getenv("EDITOR")
	}
	if editor == "" {
		editor = "vi"
	}
	cmd := exec.Command(editor, tmpPath)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("editor failed: %w", err)
	}
	data, _ := ReadFile(tmpPath)
	if !strings.HasSuffix(string(data), "\n") {
		data = append(data, '\n')
	}
	return string(data), nil
}
