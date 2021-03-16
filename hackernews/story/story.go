package story

import "github.com/theonejonahgold/pwa/hackernews"

// New creates a new Story struct
func New() hackernews.HackerNewsObject {
	return &Story{}
}

// StoryIDArray is a list of Story IDs
type StoryIDArray [500]int
