package models

import (
	"sort"
	"strconv"
	"sync"
)

// Story is a data object storing a story
type Story struct {
	ID          int        `json:"id"`
	By          string     `json:"by"`
	Descendants int        `json:"descendants"`
	Kids        []int      `json:"kids"`
	Score       int        `json:"score"`
	Time        int        `json:"time"`
	Title       string     `json:"title"`
	URL         string     `json:"url"`
	Type        string     `json:"type"`
	Deleted     bool       `json:"deleted"`
	Dead        bool       `json:"dead"`
	Comments    []*Comment `json:"comments"`
}

func (s *Story) PopulateComments() {
	kids := s.Kids
	if len(kids) == 0 {
		return
	}
	jc := make(chan string, len(kids))
	for _, v := range kids {
		jc <- strconv.Itoa(v)
	}
	close(jc)

	cc := make(chan *Comment, len(kids))
	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go fetchComment(jc, cc, &wg)
	}

	go CloseWhenWGIsDone(cc, &wg)

	cs := make([]*Comment, 0, len(kids))
	for v := range cc {
		if v.Type == "comment" {
			cs = append(cs, v)
		}
	}
	sort.Sort(CommentsByTime(cs))
	s.Comments = cs
}

func (s *Story) GetComments() []*Comment {
	return s.Comments
}

func (s *Story) GetKids() []int {
	return s.Kids
}

func (s *Story) GetType() string {
	return s.Type
}

func NewStory() *Story {
	return new(Story)
}

// StoryIDArray is a list of Story IDs
type StoryIDArray [500]int

type StoriesByTime []*Story

func (a StoriesByTime) Len() int           { return len(a) }
func (a StoriesByTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a StoriesByTime) Less(i, j int) bool { return a[i].Time > a[j].Time }

type StoriesByScore []*Story

func (a StoriesByScore) Len() int           { return len(a) }
func (a StoriesByScore) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a StoriesByScore) Less(i, j int) bool { return a[i].Score > a[j].Score }
