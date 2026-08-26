package reader

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pardnchiu/go-pkg/filesystem"
	"github.com/pardnchiu/go-pkg/utils"
)

func ListAll(dir string, opts ...ListOption) ([]File, error) {
	dir = utils.AbsPath("", dir)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("os.ReadDir: %w", err)
	}

	opt := getListOption(opts)

	out := make([]File, 0, len(entries))
	for _, e := range entries {
		if opt.SkipExcluded && filesystem.IsExcluded(dir, filepath.Join(dir, e.Name())) {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		out = append(out, newFile(filepath.Join(dir, e.Name()), info))
	}
	return out, nil
}
