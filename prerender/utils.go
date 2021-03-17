package prerender

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func createFile(path string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0770); err != nil {
		return nil, err
	}
	return os.Create(path)
}

func clearDistFolder() error {
	log.Println("Clearing dist folder")
	d, _ := os.Getwd()
	fp := filepath.Join(d, "dist")
	return os.RemoveAll(fp)
}

func saveBuildTimeToDisk() error {
	timeStamp := time.Now().Unix()
	d, _ := os.Getwd()
	fp := filepath.Join(d, "dist", "build-timestamp.txt")
	err := os.WriteFile(fp, []byte(strconv.Itoa(int(timeStamp))), 0770)
	return err
}
