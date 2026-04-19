package sandbox

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func vaildateDir(path string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("os.UserHomeDir: %w", err)
	}

	abs, err := filepath.EvalSymlinks(path)
	if err != nil {
		return "", fmt.Errorf("filepath.EvalSymlinks: %w", err)
	}

	if !strings.HasPrefix(abs, homeDir+"/") && abs != homeDir {
		return "", fmt.Errorf("just allow paths under home: %s", abs)
	}

	return homeDir, nil
}

type deniedMap struct {
	Dirs  []string `json:"dirs"`
	Files []string `json:"files"`
}

var (
	cachedDenied deniedMap
	newOnce      sync.Once
)

func New(data []byte) {
	newOnce.Do(func() {
		_ = json.Unmarshal(data, &cachedDenied)
	})
}

func deniedPaths(homeDir string) (dirs []string, files []string) {
	for _, d := range cachedDenied.Dirs {
		dirs = append(dirs, filepath.Join(homeDir, d))
	}
	for _, f := range cachedDenied.Files {
		files = append(files, filepath.Join(homeDir, f))
	}
	return
}
