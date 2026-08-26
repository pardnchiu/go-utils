package filesystem

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pardnchiu/go-pkg/utils"
)

func WriteFile(path, content string, permission os.FileMode) error {
	path = utils.AbsPath("", path)

	if IsDenied(path) {
		return fmt.Errorf("access denied: %s", path)
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("os.MkdirAll: %w", err)
	}

	// * ensure atomic write:
	// * pre-save data as temp
	tmpFile, err := os.CreateTemp(dir, "."+filepath.Base(path)+".tmp-*")
	if err != nil {
		return fmt.Errorf("os.CreateTemp: %w", err)
	}
	tmp := tmpFile.Name()
	defer os.Remove(tmp)

	if err := tmpFile.Chmod(permission); err != nil {
		tmpFile.Close()
		return fmt.Errorf("os.Chmod: %w", err)
	}
	if _, err := tmpFile.WriteString(content); err != nil {
		tmpFile.Close()
		return fmt.Errorf("os.WriteFile: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("os.Close: %w", err)
	}

	// * rename temp to target
	if err := os.Rename(tmp, path); err != nil {
		return fmt.Errorf("os.Rename: %w", err)
	}
	return nil
}
