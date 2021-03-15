package serve

import (
	"net/http"
	"os"
	"path/filepath"
)

func New() (http.Handler, error) {
	wd, _ := os.Getwd()
	fp := filepath.Join(wd, "dist")
	r := http.NewServeMux()
	r.Handle("/", http.FileServer(http.Dir(fp)))
	return r, nil
}
