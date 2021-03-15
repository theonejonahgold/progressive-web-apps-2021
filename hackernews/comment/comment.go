package comment

// CommentMap is a map with comments bound to ids
type CommentMap map[int][]Comment

type CommentsByTime []*Comment

func (a CommentsByTime) Len() int           { return len(a) }
func (a CommentsByTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a CommentsByTime) Less(i, j int) bool { return a[i].Time > a[j].Time }
