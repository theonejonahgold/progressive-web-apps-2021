package comment

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

func (c *Comment) PopulateComments(wg *sync.WaitGroup) {
	defer wg.Done()
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
	var cwg sync.WaitGroup
	for i := 0; i < 2; i++ {
		cwg.Add(1)
		go commentWorker(jc, cc, &cwg)
	}
	go func() {
		cwg.Wait()
		close(cc)
	}()
	cs := make([]*Comment, 0, len(kids))
	for v := range cc {
		if v.Type == "comment" {
			cs = append(cs, v)
		}
	}
	sort.Sort(CommentsByTime(cs))
	c.Comments = cs
	for _, v := range cs {
		wg.Add(1)
		go v.PopulateComments(wg)
	}
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

func New() *Comment {
	return &Comment{}
}
