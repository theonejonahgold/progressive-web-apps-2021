package ssg

import (
	"os"
	"path/filepath"
)

func createFile(path string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0770); err != nil {
		return nil, err
	}
	return os.Create(path)
}

func clearDistFolder() error {
	d, _ := os.Getwd()
	fp := filepath.Join(d, "dist")
	return os.RemoveAll(fp)
}
