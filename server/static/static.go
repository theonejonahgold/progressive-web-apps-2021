package static

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/NYTimes/gziphandler"
)

// New creates a new http handler for static file serving
func New() http.Handler {
	wd, _ := os.Getwd()
	fp := filepath.Join(wd, "dist")
	r := http.NewServeMux()
	r.Handle("/", gziphandler.GzipHandler(http.FileServer(http.Dir(fp))))
	r.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		bfp := filepath.Join(fp, "build-timestamp.txt")
		v, err := os.ReadFile(bfp)
		if err != nil {
			log.Println(fmt.Errorf("something went wrong while trying to fetch build timestamp: %v", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain")
		w.Write(v)
	})
	return r
}
