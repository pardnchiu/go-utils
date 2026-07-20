package filesystem

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

type AbsPathOption struct {
	HomeOnly    bool
	NeedExclude bool
}

func RealPath(path string) (string, error) {
	resolved, err := filepath.EvalSymlinks(path)
	if err == nil {
		return resolved, nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		return "", fmt.Errorf("filepath.EvalSymlinks: %w", err)
	}

	suffix := []string{filepath.Base(path)}
	dir := filepath.Dir(path)
	for {
		if dir == filepath.Dir(dir) {
			return "", fmt.Errorf("no existing ancestor for %s", path)
		}
		realAncestor, parentErr := filepath.EvalSymlinks(dir)
		if parentErr == nil {
			parts := append([]string{realAncestor}, suffix...)
			return filepath.Join(parts...), nil
		}
		if !errors.Is(parentErr, os.ErrNotExist) {
			return "", fmt.Errorf("filepath.EvalSymlinks: %w", parentErr)
		}
		suffix = append([]string{filepath.Base(dir)}, suffix...)
		dir = filepath.Dir(dir)
	}
}

func underPerUserTempRoot(resolved string) bool {
	var root string
	switch runtime.GOOS {
	case "darwin":
		root = os.TempDir()
	case "linux":
		root = os.Getenv("XDG_RUNTIME_DIR")
	}
	if root == "" {
		return false
	}
	realRoot, err := RealPath(root)
	if err != nil {
		return false
	}
	return resolved == realRoot || strings.HasPrefix(resolved, realRoot+string(filepath.Separator))
}

var wslUserRootRe = regexp.MustCompile(`^/mnt/[a-zA-Z]/Users/[^/]+(/|$)`)

func isWSLProcVersion() bool {
	raw, err := os.ReadFile("/proc/version")
	if err != nil {
		return false
	}
	return strings.Contains(strings.ToLower(string(raw)), "microsoft")
}

func underWSLWindowsUserRoot(resolved string) bool {
	if runtime.GOOS != "linux" {
		return false
	}
	if os.Getenv("WSL_DISTRO_NAME") == "" && !isWSLProcVersion() {
		return false
	}
	return wslUserRootRe.MatchString(resolved)
}

func AbsPath(root, path string, opt AbsPathOption) (string, error) {
	path = strings.TrimSpace(path)

	switch {
	case path == "":
		path = root
	case path == "~" || strings.HasPrefix(path, "~/"):
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("os.UserHomeDir: %w", err)
		}
		path = filepath.Join(home, path[1:])
	case path == "." || strings.HasPrefix(path, "./"):
		path = filepath.Join(root, strings.TrimPrefix(path, "./"))
	case !filepath.IsAbs(path):
		path = filepath.Join(root, path)
	}

	if !filepath.IsAbs(path) {
		abs, err := filepath.Abs(path)
		if err != nil {
			return "", fmt.Errorf("filepath.Abs: %w", err)
		}
		path = abs
	}

	resolved, err := RealPath(path)
	if err != nil {
		return "", fmt.Errorf("RealPath: %w", err)
	}

	if opt.HomeOnly {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("os.UserHomeDir: %w", err)
		}
		homePrefix := home + string(filepath.Separator)
		inHome := resolved == home || strings.HasPrefix(resolved, homePrefix)
		if !inHome && !underPerUserTempRoot(resolved) && !underWSLWindowsUserRoot(resolved) {
			return "", fmt.Errorf("path outside user home: %s", path)
		}
	}

	if IsDenied(resolved) {
		return "", fmt.Errorf("access denied: %s", path)
	}

	if opt.NeedExclude && IsExcluded(root, resolved) {
		return "", fmt.Errorf("path is excluded: %s", path)
	}

	return resolved, nil
}
