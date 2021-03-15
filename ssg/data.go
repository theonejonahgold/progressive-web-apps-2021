package ssg

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"sync"

	hn "github.com/theonejonahgold/pwa/hackernews"
	"github.com/theonejonahgold/pwa/hackernews/story"
)

func prepareData() (hn.DataStructure, error) {
	fmt.Println("Downloading stories")
	st, err := story.GetTopStories()
	if err != nil {
		return []*story.Story{}, err
	}

	var structure hn.DataStructure = st

	fmt.Printf("Downloading comments for %v stories (this may take a while...)\n", len(st))

	var wg sync.WaitGroup
	for _, v := range st {
		wg.Add(1)
		go v.PopulateComments(&wg)
	}
	wg.Wait()

	fmt.Println("Downloading done!")
	if err := saveDataToFile(structure); err != nil {
		fmt.Println("Data saving error:", err)
		fmt.Println("Continuing anyway")
	}
	return structure, nil
}

func saveDataToFile(structure hn.DataStructure) error {
	fmt.Println("Saving data to file.")
	j, err := json.Marshal(structure)
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
