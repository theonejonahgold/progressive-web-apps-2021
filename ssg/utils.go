package ssg

import (
	"fmt"
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

func prepareTemplate() {
	fmt.Println("Preparing Template")
	d, _ := os.Getwd()
	headFP := filepath.Join(d, "views", "partials", "head.hbs")
	headerFP := filepath.Join(d, "views", "partials", "header.hbs")
	commentFP := filepath.Join(d, "views", "partials", "comment.hbs")
	headBuf, _ := ioutil.ReadFile(headFP)
	headerBuf, _ := ioutil.ReadFile(headerFP)
	commentBuf, _ := ioutil.ReadFile(commentFP)
	raymond.RegisterPartial("partials/head", string(headBuf))
	raymond.RegisterPartial("partials/header", string(headerBuf))
	raymond.RegisterPartial("partials/comment", string(commentBuf))
}
