package keychain

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func Get(key string) string {
	if val := readKeychain(key); val != "" {
		return val
	}
	return os.Getenv(key)
}

func readKeychain(key string) string {
	switch runtime.GOOS {
	case "darwin":
		return getSecretFromMac(key)
	default:
		if secret := getSecret(key); secret != "" {
			return secret
		}
		return getFallback(key)
	}
}

func getSecretFromMac(key string) string {
	out, err := exec.Command("security", "find-generic-password",
		"-s", service,
		"-a", key,
		"-w").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func getSecret(key string) string {
	out, err := exec.Command("secret-tool", "lookup",
		"service", service, "account", key).Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func getFallback(key string) string {
	prefix := key + "="
	for _, l := range readFallbackLines() {
		if v, ok := strings.CutPrefix(l, prefix); ok {
			return v
		}
	}
	return ""
}

func readFallbackLines() []string {
	path := filepath.Join(fallbackPath, ".secrets")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var lines []string
	for line := range strings.SplitSeq(string(data), "\n") {
		if line = strings.TrimSpace(line); line != "" {
			lines = append(lines, line)
		}
	}
	return lines
}
