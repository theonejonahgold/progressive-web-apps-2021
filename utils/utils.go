package utils

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"sync"

	m "github.com/theonejonahgold/pwa/models"
)

// GetTopStories gets 50 stories from the hackernews top stories api.
func GetTopStories() ([]*m.Story, error) {
	j, err := Fetch("https://hacker-news.firebaseio.com/v0/topstories.json")
	if err != nil {
		return nil, err
	}

	var storyIDs m.StoryIDArray
	if err = json.Unmarshal(j, &storyIDs); err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	c := make(chan *m.Story, len(storyIDs))
	for _, value := range storyIDs {
		wg.Add(1)
		go topStoryFetcher(strconv.Itoa(value), c, &wg)
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	stories := make([]*m.Story, 0, len(storyIDs))
	for s := range c {
		if s == nil || s.Type != "story" {
			continue
		}
		stories = append(stories, s)
	}
	return stories, nil
}

func Fetch(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	b, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}
	return b, nil
}

func FetchComments(obj m.HackerNewsObject, wg *sync.WaitGroup, limit int) {
	defer wg.Done()

	obj.PopulateComments()
	for _, c := range obj.GetComments() {
		wg.Add(1)
		go localFetchComments(c, wg, 2, limit)
	}
}

func localFetchComments(obj m.HackerNewsObject, wg *sync.WaitGroup, level, limit int) {
	defer wg.Done()

	obj.PopulateComments()
	if limit >= 0 && limit == level {
		return
	}
	for _, c := range obj.GetComments() {
		wg.Add(1)
		go localFetchComments(c, wg, level+1, limit)
	}
}

func ParseStory(b []byte) (*m.Story, error) {
	s := m.NewStory()
	if err := json.Unmarshal(b, &s); err != nil {
		return nil, err
	}
	return s, nil
}

func topStoryFetcher(id string, c chan<- *m.Story, wg *sync.WaitGroup) {
	defer wg.Done()

	b, err := Fetch("https://hacker-news.firebaseio.com/v0/item/" + id + ".json")
	if err != nil {
		return
	}

	s, err := ParseStory(b)
	if err != nil {
		return
	}

	c <- s
}

func RetrieveSnowpackFilePath() (string, error) {
	wd, _ := os.Getwd()
	fp := filepath.Join(wd, "node_modules", ".bin", "snowpack")
	file, err := exec.LookPath(fp)
	return file, err
}
