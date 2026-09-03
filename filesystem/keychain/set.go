package keychain

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/pardnchiu/go-pkg/filesystem"
)

func Set(key, value string) error {
	if value == "" {
		return nil
	}
	switch runtime.GOOS {
	case "darwin":
		return setSecretOnMac(key, value)
	default:
		if ok := setSecret(key, value); ok == nil {
			return nil
		}
		return setFallback(key, value)
	}
}

func setSecretOnMac(key, value string) error {
	exec.Command("security", "delete-generic-password",
		"-s", service,
		"-a", key).Run()

	cmd := exec.Command("security", "add-generic-password",
		"-s", service,
		"-a", key,
		"-w", value)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("security add-generic-password: %s", out)
	}
	return nil
}

func setSecret(key, value string) error {
	ctx, cancel := context.WithTimeout(context.Background(), secretToolTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "secret-tool", "store",
		"--label", service+"/"+key,
		"service", service, "account", key)
	cmd.Stdin = strings.NewReader(value)
	cmd.WaitDelay = secretToolWaitDelay
	out, err := cmd.CombinedOutput()
	if ctx.Err() != nil {
		return fmt.Errorf("secret-tool store: %w", ctx.Err())
	}
	if err != nil {
		return fmt.Errorf("secret-tool store: %w: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

func setFallback(key, value string) error {
	path := filepath.Join(fallbackPath, ".secrets")
	lines := readFallbackLines()
	prefix := key + "="
	found := false
	for i, l := range lines {
		if strings.HasPrefix(l, prefix) {
			lines[i] = prefix + value
			found = true
			break
		}
	}
	if !found {
		lines = append(lines, prefix+value)
	}
	data := strings.Join(lines, "\n") + "\n"
	return filesystem.WriteFile(path, data, 0600)
}
