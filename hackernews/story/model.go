package story

import (
	"sort"
	"strconv"
	"sync"

	hn "github.com/theonejonahgold/pwa/hackernews"
)

// Story is a data object storing a story
type Story struct {
	ID          int                   `json:"id"`
	By          string                `json:"by"`
	Descendants int                   `json:"descendants"`
	Kids        []int                 `json:"kids"`
	Score       int                   `json:"score"`
	Time        int                   `json:"time"`
	Title       string                `json:"title"`
	URL         string                `json:"url"`
	Type        string                `json:"type"`
	Deleted     bool                  `json:"deleted"`
	Dead        bool                  `json:"dead"`
	Comments    []hn.HackerNewsObject `json:"comments"`
}

// PopulateComments
func (s *Story) PopulateComments(wg *sync.WaitGroup) {
	defer wg.Done()
	kids := s.Kids
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
	s.Comments = cs
	for _, v := range cs {
		wg.Add(1)
		go v.PopulateComments(wg)
	}
}

func (s *Story) GetType() string {
	return s.Type
}

func (s *Story) GetKids() []int {
	return s.Kids
}

func (s *Story) GetScore() int {
	return s.Score
}

func (s *Story) GetTime() int {
	return s.Time
}

func (s *Story) GetID() int {
	return s.ID
}
