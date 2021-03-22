package prerender

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	hn "github.com/theonejonahgold/pwa/hackernews"
	"github.com/theonejonahgold/pwa/renderer/handlebars"
	"github.com/theonejonahgold/pwa/snowpack"
)

var r = handlebars.NewRenderer("views")

// Build builds the entire website
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
	if _, err := favourites(); err != nil {
		return err
	}
	if err := snowpack.RunBuild(); err != nil {
		return err
	}
	if err := saveBuildTimeToDisk(); err != nil {
		return err
	}
	log.Println("Done building!")
	return nil
}

func index(data []hn.HackerNewsObject) (int, error) {
	log.Println("Rendering index page")
	bind := map[string]interface{}{
		"stories": data,
	}
	return r.Render(pageWriter{"index.html"}, "index.hbs", bind, "layouts/main.hbs")
}

func stories(data []hn.HackerNewsObject) error {
	log.Println("Rendering story pages")

	var wg sync.WaitGroup
	for _, v := range data {
		wg.Add(1)
		go storyPage(v, &wg)
	}
	wg.Wait()
	return nil
}

func storyPage(s hn.HackerNewsObject, wg *sync.WaitGroup) (int, error) {
	defer wg.Done()
	bind := map[string]interface{}{
		"story": s,
	}
	p := "story/" + strconv.Itoa(s.GetID()) + "/index.html"
	return r.Render(pageWriter{p}, "story.hbs", bind, "layouts/main.hbs")
}

func offline() (int, error) {
	return r.Render(pageWriter{"offline/index.html"}, "offline.hbs", map[string]interface{}{}, "layouts/main.hbs")
}

func favourites() (int, error) {
	return r.Render(pageWriter{"favourites/index.html"}, "favourites.hbs", map[string]interface{}{}, "layouts/main.hbs")
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
