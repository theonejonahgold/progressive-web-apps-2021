package serve

import (
	"encoding/json"
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
	r.HandleFunc("/all-stories", allStoriesHandler)
	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:" + port,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
	}
	return srv.ListenAndServe()
}

func allStoriesHandler(w http.ResponseWriter, req *http.Request) {
	d, _ := os.Getwd()
	fp := filepath.Join(d, "dist", "story")
	file, err := os.Open(fp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong"))
		return
	}
	defer file.Close()
	list, _ := file.Readdirnames(0)
	for k := range list {
		list[k] = "/story/" + list[k] + "/"
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(w)
	enc.Encode(list)
}
