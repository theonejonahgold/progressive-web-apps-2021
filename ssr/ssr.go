package ssr

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"github.com/theonejonahgold/pwa/hackernews/story"
	"github.com/theonejonahgold/pwa/renderer"
	"github.com/theonejonahgold/pwa/snowpack"
)

var r = renderer.New("views")

func New(ctx context.Context) (http.Handler, error) {
	err := snowpack.RunDev(ctx)
	if err != nil {
		return nil, err
	}
	r := http.NewServeMux()
	r.HandleFunc("/offline", offlineHandler)
	r.HandleFunc("/story/", storyHandler)
	r.HandleFunc("/", indexHandler)
	r.HandleFunc("/serviceWorker.js", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNotFound) })
	return r, nil
}

func indexHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		wd, _ := os.Getwd()
		fp := filepath.Join(wd, "dist")
		http.FileServer(http.Dir(fp)).ServeHTTP(w, req)
		return
	}
	stories, err := story.GetTopStories()
	if err != nil {
		return
	}
	sort.Sort(story.ByScore(stories))

	if err != nil {
		return
	}
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
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

	if err != nil {
		return
	}
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	r.Render(w, "story.hbs", map[string]interface{}{
		"story": story,
	}, "layouts/main.hbs")
}

func offlineHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	r.Render(w, "offline.hbs", map[string]interface{}{}, "layouts/main.hbs")
}

func notFoundHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(404)
	w.Write([]byte("Kon die stoorie niet vinden. Probeer eens een andere!"))
}
