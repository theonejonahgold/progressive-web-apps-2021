package story

import (
	"sort"
	"strconv"
	"sync"

	"github.com/theonejonahgold/pwa/hackernews/comment"
)

// Story is a data object storing a story
type Story struct {
	ID          int                `json:"id"`
	By          string             `json:"by"`
	Descendants int                `json:"descendants"`
	Kids        []int              `json:"kids"`
	Score       int                `json:"score"`
	Time        int                `json:"time"`
	Title       string             `json:"title"`
	URL         string             `json:"url"`
	Type        string             `json:"type"`
	Deleted     bool               `json:"deleted"`
	Dead        bool               `json:"dead"`
	Comments    []*comment.Comment `json:"comments"`
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
	cc := make(chan *comment.Comment, len(kids))
	var cwg sync.WaitGroup
	for i := 0; i < 2; i++ {
		cwg.Add(1)
		go commentWorker(jc, cc, &cwg)
	}
	go func() {
		cwg.Wait()
		close(cc)
	}()
	cs := make([]*comment.Comment, 0, len(kids))
	for v := range cc {
		if v.Type == "comment" {
			cs = append(cs, v)
		}
	}
	sort.Sort(comment.CommentsByTime(cs))
	s.Comments = cs
	for _, v := range cs {
		wg.Add(1)
		go v.PopulateComments(wg)
	}
}

func (s *Story) GetComments() []*comment.Comment {
	return s.Comments
}

func (s *Story) GetKids() []int {
	return s.Kids
}

func (s *Story) GetType() string {
	return s.Type
}
