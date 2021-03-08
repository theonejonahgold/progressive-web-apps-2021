package utils

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	m "github.com/theonejonahgold/pwa/models"
)

// RetrieveCommentsForStory retrieves all comments for a given Story
func RetrieveCommentsForStory(s *m.Story) (*[]m.Comment, error) {
	c := make(chan *m.Comment, len(s.Kids))
	for _, value := range s.Kids {
		go fetchComment(strconv.Itoa(value), c)
	}

	var comments []m.Comment
	for cmt := range c {
		comments = append(comments, *cmt)
		if len(comments) == len(s.Kids) {
			close(c)
		}
	}

	return &comments, nil
}

// GetTopStories gets 50 stories from the hackernews top stories api.
func GetTopStories() (*[]m.Story, error) {
	res, err := http.Get("https://hacker-news.firebaseio.com/v0/topstories.json")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var storyIDs m.StoryIDArray
	j, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(j, &storyIDs)
	if err != nil {
		return nil, err
	}

	c := make(chan *m.Story, 500)
	for _, value := range storyIDs {
		go fetchStory(strconv.Itoa(value), c)
	}

	var stories []m.Story
	for s := range c {
		stories = append(stories, *s)
		if len(stories) == len(storyIDs) {
			close(c)
		}
	}
	return &stories, nil
}

func fetchStory(id string, c chan *m.Story) error {
	res, err := http.Get("https://hacker-news.firebaseio.com/v0/item/" + id + ".json")
	if err != nil {
		return err
	}
	defer res.Body.Close()

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var s m.Story
	err = json.Unmarshal(bytes, &s)
	if err != nil {
		return err
	}
	c <- &s
	return nil
}

func fetchComment(id string, c chan *m.Comment) error {
	res, err := http.Get("https://hacker-news.firebaseio.com/v0/item/" + id + ".json")
	if err != nil {
		return err
	}
	defer res.Body.Close()

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var comment m.Comment
	err = json.Unmarshal(bytes, &comment)
	if err != nil {
		return err
	}
	c <- &comment
	return nil
}
