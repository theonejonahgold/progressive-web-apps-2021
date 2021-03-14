package models

import (
	"sort"
	"strconv"
	"sync"
)

// Comment - A Hacker News comment struct
type Comment struct {
	ID       int        `json:"id"`
	By       string     `json:"by"`
	Type     string     `json:"type"`
	Time     int        `json:"time"`
	Parent   int        `json:"parent"`
	Text     string     `json:"text"`
	Kids     []int      `json:"kids"`
	Deleted  bool       `json:"deleted"`
	Dead     bool       `json:"dead"`
	Comments []*Comment `json:"comments"`
}

func (c *Comment) PopulateComments() {
	kids := c.Kids
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
	c.Comments = cs
}

func (c *Comment) GetComments() []*Comment {
	return c.Comments
}

func (c *Comment) GetKids() []int {
	return c.Kids
}

func (c *Comment) GetType() string {
	return c.Type
}

func NewComment() *Comment {
	return &Comment{}
}

// CommentMap is a map with comments bound to ids
type CommentMap map[int][]Comment

type CommentsByTime []*Comment

func (a CommentsByTime) Len() int           { return len(a) }
func (a CommentsByTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a CommentsByTime) Less(i, j int) bool { return a[i].Time > a[j].Time }
