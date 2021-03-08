package models

// Story is a data object storing a story
type Story struct {
	ID          int    `json:"id"`
	By          string `json:"by"`
	Descendants int    `json:"descendants"`
	Kids        []int  `json:"kids"`
	Score       int    `json:"score"`
	Time        int    `json:"time"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	Type        string `json:"type"`
	Deleted     bool   `json:"deleted"`
	Dead        bool   `json:"dead"`
}

// StoryIDArray is a list of Story IDs
type StoryIDArray []int
