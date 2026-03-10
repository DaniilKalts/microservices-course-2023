package env

import (
	"errors"
	"fmt"
	"os"
	"time"
)

func ReadDuration(envName string) (time.Duration, error) {
	raw := os.Getenv(envName)
	if len(raw) == 0 {
		return 0, errors.New(envName + " is not set")
	}

	value, err := time.ParseDuration(raw)
	if err != nil {
		return 0, fmt.Errorf("%s must be a valid duration, got %q: %w", envName, raw, err)
	}

	if value < 0 {
		return 0, fmt.Errorf("%s must be non-negative", envName)
	}

	return value, nil
}
