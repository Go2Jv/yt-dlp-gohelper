package history

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

func FindTempFiles(outputDir string, since time.Time) ([]string, error) {
	entries, err := os.ReadDir(outputDir)
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, 16)
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !looksLikeTemp(name) {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		if !since.IsZero() && info.ModTime().Before(since.Add(-10*time.Minute)) {
			continue
		}
		out = append(out, filepath.Join(outputDir, name))
		if len(out) >= 200 {
			break
		}
	}
	return out, nil
}

func DeleteFiles(files []string) {
	for _, f := range files {
		_ = os.Remove(f)
	}
}

func AnyExists(files []string) bool {
	for _, f := range files {
		if _, err := os.Stat(f); err == nil {
			return true
		}
	}
	return false
}

func looksLikeTemp(name string) bool {
	if strings.Contains(name, ".part") {
		return true
	}
	switch {
	case strings.HasSuffix(name, ".ytdl"):
		return true
	case strings.HasSuffix(name, ".aria2"):
		return true
	case strings.HasSuffix(name, ".temp"):
		return true
	default:
		return false
	}
}

