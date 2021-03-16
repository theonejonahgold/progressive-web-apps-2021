package static

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// New creates a new http handler for static file serving
func New() http.Handler {
	wd, _ := os.Getwd()
	fp := filepath.Join(wd, "dist")
	r := http.NewServeMux()
	r.Handle("/", http.FileServer(http.Dir(fp)))
	r.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		v := strconv.Itoa(time.Now().Hour()) + "-" + strconv.Itoa(time.Now().Day())
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(v))
	})
	return r
}
