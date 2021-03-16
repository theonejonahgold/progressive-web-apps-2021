package prerender

import (
	"encoding/json"
	"log"
	"path/filepath"
	"sync"

	hn "github.com/theonejonahgold/pwa/hackernews"
	"github.com/theonejonahgold/pwa/hackernews/story"
)

func prepareData() ([]hn.HackerNewsObject, error) {
	log.Println("Downloading stories")
	st, err := story.GetTopStories()
	if err != nil {
		return []hn.HackerNewsObject{}, err
	}

	log.Printf("Downloading comments for %v stories (this may take a while...)\n", len(st))

	var wg sync.WaitGroup
	for _, v := range st {
		wg.Add(1)
		go v.PopulateComments(&wg)
	}
	wg.Wait()

	log.Println("Downloading done")
	if err := saveDataToFile(st); err != nil {
		log.Printf("data saving error: %v", err)
		log.Println("Continuing anyway")
	}
	return st, nil
}

func saveDataToFile(stories []hn.HackerNewsObject) error {
	log.Println("Saving data to file")
	j, err := json.Marshal(stories)
	if err != nil {
		return err
	}

	file, err := createFile(filepath.Join("dist", "data.json"))
	if err != nil {
		return err
	}

	if _, err = file.Write(j); err != nil {
		return err
	}
	return nil
}
