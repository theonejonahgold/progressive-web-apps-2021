package hackernews

import (
	"sync"
)

type HackerNewsObject interface {
	PopulateComments(*sync.WaitGroup)
	GetID() int
	GetType() string
	GetKids() []int
	GetTime() int
	GetScore() int
}

type ByTime []HackerNewsObject

func (a ByTime) Len() int           { return len(a) }
func (a ByTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByTime) Less(i, j int) bool { return a[i].GetTime() > a[j].GetTime() }

type ByScore []HackerNewsObject

func (a ByScore) Len() int           { return len(a) }
func (a ByScore) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByScore) Less(i, j int) bool { return a[i].GetScore() > a[j].GetScore() }
