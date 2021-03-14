package ssg

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"sync"

	m "github.com/theonejonahgold/pwa/models"
	u "github.com/theonejonahgold/pwa/utils"
)

func prepareData() (m.DataStructure, error) {
	fmt.Println("Downloading data")
	st, err := u.GetTopStories()
	if err != nil {
		return []*m.Story{}, err
	}
	fmt.Println("Stories downloaded:", len(st))

	var structure m.DataStructure = st

	var wg sync.WaitGroup
	for _, v := range st {
		wg.Add(1)
		go u.FetchComments(v, &wg, -1)
	}

	wg.Wait()

	fmt.Println("Downloading done")

	if err := saveDumpToFile(structure); err != nil {
		return []*m.Story{}, err
	}
	return structure, nil
}

func saveDumpToFile(structure m.DataStructure) error {
	fmt.Println("Saving data")
	j, err := json.Marshal(structure)
	if err != nil {
		return err
	}

	file, err := createFile(filepath.Join("dist", "data-dump.json"))
	if err != nil {
		return err
	}

	if _, err = file.Write(j); err != nil {
		return err
	}
	return nil
}
