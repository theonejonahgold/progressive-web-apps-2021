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
	fmt.Println("Downloading stories")
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

	fmt.Println("Downloading comments (this may take a while...)")
	wg.Wait()

	fmt.Println("Downloading done!")
	if err := saveDataToFile(structure); err != nil {
		fmt.Println("Data saving error:", err)
		fmt.Println("Continuing anyway")
	}
	return structure, nil
}

func saveDataToFile(structure m.DataStructure) error {
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
