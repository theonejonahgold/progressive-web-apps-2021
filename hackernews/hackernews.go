package hackernews

import (
	"sync"

	"github.com/theonejonahgold/pwa/hackernews/comment"
)

type HackerNewsObject interface {
	PopulateComments(*sync.WaitGroup)
	GetComments() []*comment.Comment
	GetKids() []int
	GetType() string
}
