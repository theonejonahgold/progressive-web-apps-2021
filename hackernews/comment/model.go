package comment

import (
	"sort"
	"strconv"
	"sync"

	hn "github.com/theonejonahgold/pwa/hackernews"
)

// Comment - A Hacker News comment struct
type Comment struct {
	ID       int                   `json:"id"`
	By       string                `json:"by"`
	Type     string                `json:"type"`
	Time     int                   `json:"time"`
	Parent   int                   `json:"parent"`
	Text     string                `json:"text"`
	Kids     []int                 `json:"kids"`
	Deleted  bool                  `json:"deleted"`
	Dead     bool                  `json:"dead"`
	Comments []hn.HackerNewsObject `json:"comments"`
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
	cc := make(chan hn.HackerNewsObject, len(kids))
	var cwg sync.WaitGroup
	for i := 0; i < len(kids); i++ {
		cwg.Add(1)
		go commentWorker(jc, cc, &cwg)
	}
	go func() {
		cwg.Wait()
		close(cc)
	}()
	cs := make([]hn.HackerNewsObject, 0, len(kids))
	for v := range cc {
		if v.GetType() == "comment" {
			cs = append(cs, v)
		}
	}
	sort.Sort(hn.ByTime(cs))
	c.Comments = cs
	for _, v := range cs {
		wg.Add(1)
		go v.PopulateComments(wg)
	}
}

func (c *Comment) GetType() string {
	return c.Type
}

func (c *Comment) GetKids() []int {
	return c.Kids
}

func (c *Comment) GetTime() int {
	return c.Time
}

func (c *Comment) GetScore() int {
	return 0
}

func (c *Comment) GetID() int {
	return c.ID
}

func New() hn.HackerNewsObject {
	return &Comment{}
}
