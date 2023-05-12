package rest

type CreateItem struct {
	Text      string
	State     string
	SortIndex int
}

type UpdateItem struct {
	Text      string
	State     string
	SortIndex int
}
