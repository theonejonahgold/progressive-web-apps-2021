package models

type HackerNewsObject interface {
	PopulateComments()
	GetComments() []*Comment
	GetKids() []int
	GetType() string
}
