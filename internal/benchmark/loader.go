package benchmark

import (
	"os"
	"path/filepath"
	"strings"
)

func LoadSuite(dir string) []string {
	var files []string

	entries, err := os.ReadDir(dir)
	if err != nil {
		_ = os.MkdirAll(dir, os.ModePerm)
		return []string{}
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			name := entry.Name()
			if strings.HasSuffix(name, ".max") {
				files = append(files, filepath.Join(dir, name))
			}
		}
	}
	return files
}
