package serve

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func Serve() error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	fmt.Println("Starting serve on port", port)
	wd, _ := os.Getwd()
	fp := filepath.Join(wd, "dist")
	r := http.NewServeMux()
	r.Handle("/", http.FileServer(http.Dir(fp)))
	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:" + port,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
	}
	return srv.ListenAndServe()
}
