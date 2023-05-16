package rest

import (
	"time"

	"github.com/tagptroll1/groupie-api/model/dbmodel"
)

type Item struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	ListID    string    `json:"listId"`
	Text      string    `json:"text"`
	State     string    `json:"state"`
	SortIndex int       `json:"sortIndex"`
}

func ToItem(i dbmodel.Item) Item {
	return Item{
		ID:        i.ID,
		CreatedAt: i.CreatedAt,
		UpdatedAt: i.UpdatedAt,
		ListID:    i.ListID,
		Text:      i.Text,
		State:     i.State,
		SortIndex: i.SortIndex,
	}
}

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
