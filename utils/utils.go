package utils

import (
	"path/filepath"
	"strings"
)

func ToFilePath(str string) string {
	return filepath.FromSlash(strings.ReplaceAll(str, `\`, `/`))
}
