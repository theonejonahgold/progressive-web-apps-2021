package models

type HackerNewsObject interface {
	PopulateComments()
	GetKids() []int
	GetType() string
}
