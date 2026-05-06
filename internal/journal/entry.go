package journal

import (
	"os"
	"path/filepath"
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

