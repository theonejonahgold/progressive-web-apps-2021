package utils

import (
	"encoding/json"
	"io"
	"net/http"
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
		go fetchStory(strconv.Itoa(value), c, &wg)
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	var stories []*m.Story
	for s := range c {
		if s.Type == "story" {
			stories = append(stories, s)
		}
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

func FetchComments(obj m.HackerNewsObject, wg *sync.WaitGroup) {
	defer wg.Done()
	obj.PopulateComments()
	switch v := (obj).(type) {
	case *m.Story:
		for _, c := range v.Comments {
			wg.Add(1)
			go FetchComments(c, wg)
		}
	case *m.Comment:
		for _, c := range v.Comments {
			wg.Add(1)
			go FetchComments(c, wg)
		}
	}
}

func ParseStory(b []byte) (*m.Story, error) {
	s := m.NewStory()
	if err := json.Unmarshal(b, &s); err != nil {
		return s, err
	}
	return s, nil
}

func ParseComment(b []byte) (*m.Comment, error) {
	c := m.NewComment()
	if err := json.Unmarshal(b, &c); err != nil {
		return c, err
	}
	return c, nil
}

func fetchStory(id string, c chan *m.Story, wg *sync.WaitGroup) error {
	defer wg.Done()

	b, err := Fetch("https://hacker-news.firebaseio.com/v0/item/" + id + ".json")
	if err != nil {
		return err
	}

	s, err := ParseStory(b)
	if err != nil {
		return err
	}

	c <- s
	return nil
}
