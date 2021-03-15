package utils

import (
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

func Fetch(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	b, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}
	return b, nil
}

func RetrieveSnowpackFilePath() (string, error) {
	wd, _ := os.Getwd()
	fp := filepath.Join(wd, "node_modules", ".bin", "snowpack")
	file, err := exec.LookPath(fp)
	return file, err
}
