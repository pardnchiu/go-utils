package filesystem

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pardnchiu/go-pkg/utils"
)

func ReadJSON[T any](path string) (T, error) {
	var v T
	bytes, err := os.ReadFile(utils.AbsPath("", path))
	if err != nil {
		return v, fmt.Errorf("os.ReadFile: %w", err)
	}

	if err := json.Unmarshal(bytes, &v); err != nil {
		return v, fmt.Errorf("json.Unmarshal: %w", err)
	}
	return v, nil
}

func WriteJSON(path string, v any, format bool) error {
	var (
		bytes []byte
		err   error
	)
	if format {
		bytes, err = json.MarshalIndent(v, "", "  ")
	} else {
		bytes, err = json.Marshal(v)
	}
	if err != nil {
		return fmt.Errorf("json.Marshal: %w", err)
	}
	return WriteFile(path, string(bytes), 0644)
}
