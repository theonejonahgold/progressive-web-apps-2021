package static

import (
	"net/http"
	"os"
	"path/filepath"
)

// New creates a new http handler for static file serving
func New() http.Handler {
	wd, _ := os.Getwd()
	fp := filepath.Join(wd, "dist")
	r := http.NewServeMux()
	r.Handle("/", http.FileServer(http.Dir(fp)))
	return r
}
