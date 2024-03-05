package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/mattn/go-isatty"
)

const (
	ExitOk int = iota
	ExitErr
	ExitCancel
)

// A map for referencing values by env var if they exist
var defaultMap = map[string]string{
	"authToken":        "GH_AUTH_TOKEN",
	"defaultBaseURL":   "GH_BASE_URL",
	"defaultUserAgent": "GH_USER_AGENT",
}

// For logging purposes, GetEnvIfSet returns the env var name and its value
func GetEnvIfSet(key string) (name string, val string, err error) {
	envVar, ok := defaultMap[key]
	if !ok {
		// This is to catch programmer error, like a typo in a key name
		return "", "", fmt.Errorf("referenced unknown config key: %s", key)
	}

	return envVar, os.Getenv(envVar), nil
}

func TTYVerify() error {
	if isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd()) {
		return nil
	}
	return errors.New("compatible terminal environment not found")
}
