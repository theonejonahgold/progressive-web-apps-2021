package ssr

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/theonejonahgold/pwa/hackernews/story"
	"github.com/theonejonahgold/pwa/renderer"
	"github.com/theonejonahgold/pwa/snowpack"
)

func SSR() (context.Context, error) {
	ctx := context.Background()
	err := snowpack.RunDev(ctx)
	if err != nil {
		return nil, err
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	fmt.Println("Starting serve on port", port)

	wd, _ := os.Getwd()
	fp := filepath.Join(wd, "dist", "assets")
	r := http.NewServeMux()
	r.HandleFunc("/story/", storyHandler)
	r.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir(fp))))
	r.HandleFunc("/", indexHandler)
	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:" + port,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
	}
	return ctx, srv.ListenAndServe()
}

func indexHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		notFoundHandler(w, req)
		return
	}
	stories, err := story.GetTopStories()
	if err != nil {
		return
	}
	sort.Sort(story.ByScore(stories))
	r, err := renderer.New("views")
	if err != nil {
		return
	}
	w.Header().Add("Content-Type", "text/html")
	r.Render(w, "index.hbs", map[string]interface{}{
		"stories": stories,
	}, "layouts/main.hbs")
}

func storyHandler(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Path[7:]
	if id == "" {
		notFoundHandler(w, req)
		return
	}
	res, err := http.Get("https://hacker-news.firebaseio.com/v0/item/" + id + ".json")
	if err != nil {
		return
	}

	story, err := story.Parse(res)
	if err != nil {
		return
	}

	if story.Type != "story" {
		notFoundHandler(w, req)
		return
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go story.PopulateComments(&wg)
	wg.Wait()

	r, err := renderer.New("views")
	if err != nil {
		return
	}
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	r.Render(w, "story.hbs", map[string]interface{}{
		"story": story,
	}, "layouts/main.hbs")
}

func notFoundHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(404)
	w.Write([]byte("Kon die stoorie niet vinden. Probeer eens een andere!"))
}
