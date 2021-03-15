package ssg

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/theonejonahgold/pwa/hackernews/story"
	"github.com/theonejonahgold/pwa/renderer"
	"github.com/theonejonahgold/pwa/snowpack"
)

// GeneratePages generates all pages for static rendering
func SSG() error {
	data, err := prepareData()
	if err != nil {
		return err
	}
	if _, err := renderIndex(data); err != nil {
		return err
	}
	if err := renderStories(data); err != nil {
		return err
	}
	if err := snowpack.RunBuild(); err != nil {
		return err
	}
	fmt.Println("Done building!")
	return nil
}

func renderIndex(data []*story.Story) (int, error) {
	fmt.Println("Rendering Index Page")
	r, err := renderer.New("views")
	if err != nil {
		return 0, err
	}
	bind := map[string]interface{}{
		"stories": data,
	}
	return r.Render(fileSaver{
		path: "index.html",
	}, "index.hbs", bind, "layouts/main.hbs")
}

func renderStories(data []*story.Story) error {
	fmt.Println("Rendering Story Pages")
	r, err := renderer.New("views")
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for _, v := range data {
		wg.Add(1)
		go renderStory(r, v, &wg)
	}
	wg.Wait()
	return nil
}

func renderStory(r renderer.Renderer, s *story.Story, wg *sync.WaitGroup) (int, error) {
	defer wg.Done()
	bind := map[string]interface{}{
		"story": s,
	}
	return r.Render(fileSaver{
		path: "story/" + strconv.Itoa(s.ID) + "/index.html",
	}, "story.hbs", bind, "layouts/main.hbs")
}

type fileSaver struct {
	path string
}

func (s fileSaver) Write(data []byte) (n int, err error) {
	d, _ := os.Getwd()
	fp := filepath.Join(d, "dist", s.path)
	f, err := createFile(fp)
	if err != nil {
		return 0, err
	}
	n, err = f.Write(data)
	if err != nil {
		return 0, err
	}
	if err = f.Close(); err != nil {
		return 0, err
	}
	return n, nil
}
