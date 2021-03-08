package models

// Comment is a data object storing a comment
type Comment struct {
	ID      int    `json:"id"`
	By      string `json:"by"`
	Type    string `json:"type"`
	Time    int    `json:"time"`
	Parent  int    `json:"parent"`
	Text    string `json:"text"`
	Kids    []int  `json:"kids"`
	Deleted bool   `json:"deleted"`
	Dead    bool   `json:"dead"`
}
