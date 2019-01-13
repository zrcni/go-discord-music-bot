package utils

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// GetBasePath returns the path where the program was executed
func GetBasePath() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", errors.Wrap(err, "getBasePath")
	}
	path := filepath.Dir(ex)

	return path, nil
}
