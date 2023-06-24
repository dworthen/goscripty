package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ToFilePath(str string) string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Unable to load user home dir.")
		os.Exit(1)
	}
	return filepath.Clean(strings.ReplaceAll(str, `~`, homeDir))
}
