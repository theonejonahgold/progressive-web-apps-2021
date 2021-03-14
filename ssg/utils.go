package ssg

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/aymerick/raymond"
)

func createFile(path string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0770); err != nil {
		return nil, err
	}
	return os.Create(path)
}

func prepareTemplate() (*raymond.Template, error) {
	fmt.Println("Preparing Template")
	d, _ := os.Getwd()
	partialP := filepath.Join(d, "views", "partials")
	filepath.Walk(partialP, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		base := filepath.Base(partialP)
		name := base + "/" + filepath.Base(path[:len(path)-4])
		buf, _ := ioutil.ReadFile(path)
		raymond.RegisterPartial(name, string(buf))
		return nil
	})
	mainLay := filepath.Join(d, "views", "layouts", "main.hbs")

	return raymond.ParseFile(mainLay)
}
