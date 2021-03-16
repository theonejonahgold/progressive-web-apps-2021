package handlebars

import (
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/aymerick/raymond"
	"github.com/theonejonahgold/pwa/renderer"
)

func NewRenderer(basePath string) (r renderer.Renderer) {
	d, _ := os.Getwd()
	partialP := filepath.Join(d, basePath, "partials")
	filepath.Walk(partialP, func(path string, info fs.FileInfo, err error) error {
		defer recoverFromPartialRegisteredError()
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
	r = &handlebarsRenderer{
		basePath,
	}
	return r
}

type handlebarsRenderer struct {
	basePath string
}

func (hr *handlebarsRenderer) Render(output io.Writer, fp string, data map[string]interface{}, layout string) (int, error) {
	template, err := raymond.ParseFile(filepath.Join(hr.basePath, fp))
	if err != nil {
		return 0, err
	}
	embed, err := template.Exec(data)
	if err != nil {
		return 0, err
	}
	data["embed"] = raymond.SafeString(embed)
	lay, err := raymond.ParseFile(filepath.Join(hr.basePath, layout))
	if err != nil {
		return 0, err
	}
	render, err := lay.Exec(data)
	if err != nil {
		return 0, err
	}
	return output.Write([]byte(render))
}

func recoverFromPartialRegisteredError() {
	recover()
}
