package ssg

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/aymerick/raymond"
	m "github.com/theonejonahgold/pwa/models"
	u "github.com/theonejonahgold/pwa/utils"
)

// GeneratePages generates all pages for static rendering
func SSG() error {
	data, err := prepareData()
	if err != nil {
		return err
	}
	main, err := prepareTemplate()
	if err != nil {
		return err
	}
	if err := renderIndex(main, data); err != nil {
		return err
	}
	if err := renderStories(main, data); err != nil {
		return err
	}
	if err := executeSnowpackBuild(); err != nil {
		return err
	}
	fmt.Println("Site ready for deployment!")
	return nil
}

func renderIndex(lay *raymond.Template, data []*m.Story) error {
	fmt.Println("Rendering Index Page")
	d, _ := os.Getwd()
	t, err := raymond.ParseFile(filepath.Join(d, "views", "index.hbs"))
	if err != nil {
		return err
	}

	bind := map[string]interface{}{
		"stories": data,
	}
	e, err := t.Exec(bind)
	if err != nil {
		return err
	}

	bind["embed"] = raymond.SafeString(e)
	r, err := lay.Exec(bind)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return savePage([]byte(r), "index.html")
}

func renderStories(lay *raymond.Template, data []*m.Story) error {
	fmt.Println("Rendering Story Pages")
	d, _ := os.Getwd()
	t, err := raymond.ParseFile(filepath.Join(d, "views", "story.hbs"))
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for _, v := range data {
		wg.Add(1)
		go renderStory(lay, t, v, &wg)
	}
	wg.Wait()
	return nil
}

func renderStory(lay *raymond.Template, t *raymond.Template, s *m.Story, wg *sync.WaitGroup) error {
	defer wg.Done()

	bind := map[string]interface{}{
		"story": s,
	}
	e, err := t.Exec(bind)
	if err != nil {
		return err
	}

	bind["embed"] = raymond.SafeString(e)
	r, err := lay.Exec(bind)
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
	file, err := u.RetrieveSnowpackFilePath()
	if err != nil {
		return err
	}
	cmd := exec.Command(file, "build")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), "NODE_ENV=production")
	return cmd.Run()
}
