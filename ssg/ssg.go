package ssg

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	s "github.com/theonejonahgold/pwa/hackernews/story"
	"github.com/theonejonahgold/pwa/renderer"
	"github.com/theonejonahgold/pwa/snowpack"
)

var r = renderer.New("views")

// Build generates all pages for static rendering
func Build() error {
	if err := clearDistFolder(); err != nil {
		return err
	}
	data, err := prepareData()
	if err != nil {
		return err
	}
	if _, err := index(data); err != nil {
		return err
	}
	if err := stories(data); err != nil {
		return err
	}
	if _, err := offline(); err != nil {
		return err
	}
	if err := snowpack.RunBuild(); err != nil {
		return err
	}
	fmt.Println("Done building!")
	return nil
}

func index(data []*s.Story) (int, error) {
	fmt.Println("Rendering Index Page")
	bind := map[string]interface{}{
		"stories": data,
	}
	return r.Render(pageWriter{
		path: "index.html",
	}, "index.hbs", bind, "layouts/main.hbs")
}

func stories(data []*s.Story) error {
	fmt.Println("Rendering Story Pages")

	var wg sync.WaitGroup
	for _, v := range data {
		wg.Add(1)
		go storyPage(v, &wg)
	}
	wg.Wait()
	return nil
}

func storyPage(s *s.Story, wg *sync.WaitGroup) (int, error) {
	defer wg.Done()
	bind := map[string]interface{}{
		"story": s,
	}
	return r.Render(pageWriter{
		path: "story/" + strconv.Itoa(s.ID) + "/index.html",
	}, "story.hbs", bind, "layouts/main.hbs")
}

func offline() (int, error) {
	return r.Render(pageWriter{
		path: "offline/index.html",
	}, "offline.hbs", map[string]interface{}{}, "layouts/main.hbs")
}

type pageWriter struct {
	path string
}

func (s pageWriter) Write(data []byte) (n int, err error) {
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
