package main

import (
	"fmt"
	"os"
	"os/exec"
)

// openEditor opens the specified file in the user's preferred editor
func openEditor(filePath, editorOverride string) error {
	editor := editorOverride
	if editor == "" {
		editor = os.Getenv("EDITOR")
	}
	if editor == "" {
		editor = "vi"
	}

	cmd := exec.Command(editor, filePath)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("editor failed: %w", err)
	}
	return nil
}
