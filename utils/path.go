package utils

import (
	"os"
	"path/filepath"
	"strings"
)

func AbsPath(base, path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		path = strings.TrimSpace(base)
		base = ""
	}
	if path == "" {
		return ""
	}

	switch {
	case filepath.IsAbs(path):
		return filepath.Clean(path)
	case path == "~" || strings.HasPrefix(path, "~/"):
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(home, path[1:])
	case base != "":
		return filepath.Join(AbsPath("", base), path)
	}

	abs, err := filepath.Abs(path)
	if err != nil {
		return filepath.Clean(path)
	}
	return abs
}
