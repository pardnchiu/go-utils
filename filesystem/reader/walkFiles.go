package reader

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pardnchiu/go-pkg/filesystem"
	"github.com/pardnchiu/go-pkg/utils"
)

func WalkFiles(root string, opts ...ListOption) ([]File, error) {
	root = utils.AbsPath("", root)
	opt := getListOption(opts)

	var files []File
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			if opt.IgnoreWalkError {
				return nil
			}
			return err
		}
		if path == root {
			return nil
		}
		if entry.IsDir() {
			if strings.HasPrefix(entry.Name(), ".") {
				return filepath.SkipDir
			}
			if opt.SkipDenied && filesystem.IsDenied(path) {
				return filepath.SkipDir
			}
			if opt.SkipExcluded && filesystem.IsExcluded(root, path) {
				return filepath.SkipDir
			}
			return nil
		}
		if !opt.IncludeNonRegular && !entry.Type().IsRegular() {
			return nil
		}
		if opt.SkipExcluded && filesystem.IsExcluded(root, path) {
			return nil
		}

		info, err := entry.Info()
		if err != nil {
			if opt.IgnoreWalkError {
				return nil
			}
			return fmt.Errorf("entry.Info: %w", err)
		}
		files = append(files, newFile(path, info))

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("filepath.WalkDir: %w", err)
	}
	return files, nil
}
