package ssg

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/aymerick/raymond"
	m "github.com/theonejonahgold/pwa/models"
)

// GeneratePages generates all pages for static rendering
func SSG() error {
	data, err := prepareData()
	if err != nil {
		return err
	}
	prepareTemplate()
	if err := renderIndex(data); err != nil {
		return err
	}
	if err := renderStories(data); err != nil {
		return err
	}
	if err := executeSnowpackBuild(); err != nil {
		return err
	}
	fmt.Println("Site ready for deployment!")
	return nil
}

func renderIndex(data []*m.Story) error {
	fmt.Println("Rendering Index Page")
	d, _ := os.Getwd()
	t, err := ioutil.ReadFile(filepath.Join(d, "views", "index.hbs"))
	if err != nil {
		return err
	}
	r, err := raymond.Render(string(t), map[string]interface{}{
		"stories": data,
	})
	if err != nil {
		return err
	}
	return savePage([]byte(r), "index.html")
}

func renderStories(data []*m.Story) error {
	fmt.Println("Rendering Story Pages")
	d, _ := os.Getwd()
	t, err := ioutil.ReadFile(filepath.Join(d, "views", "story.hbs"))
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for _, v := range data {
		wg.Add(1)
		go renderStory(v, t, &wg)
	}
	wg.Wait()
	return nil
}

func renderStory(s *m.Story, t []byte, wg *sync.WaitGroup) error {
	defer wg.Done()

	r, err := raymond.Render(string(t), map[string]interface{}{
		"story": s,
	})
	if err != nil {
		fmt.Println(err)
		return err
	}
	return savePage([]byte(r), "story/"+strconv.Itoa(s.ID)+"/index.html")
}

func savePage(t []byte, p string) error {
	d, _ := os.Getwd()
	fp := filepath.Join(d, "dist", p)
	f, err := createFile(fp)
	if err != nil {
		return err
	}
	if _, err = f.Write(t); err != nil {
		return err
	}
	if err = f.Close(); err != nil {
		return err
	}
	return err
}

func executeSnowpackBuild() error {
	fmt.Println("Running Snowpack Build")
	npmp, err := exec.LookPath("npm")
	if err != nil {
		return err
	}
	wd, _ := os.Getwd()
	cmd := exec.Cmd{
		Path:   npmp,
		Args:   []string{"npm", "run", "build:snowpack"},
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Dir:    wd,
	}
	return cmd.Run()
}
