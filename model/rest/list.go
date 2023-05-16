package rest

import (
	"time"

	"github.com/tagptroll1/groupie-api/model/dbmodel"
)

type List struct {
	ID        string           `json:"id"`
	Title     string           `json:"title"`
	Type      dbmodel.ListType `json:"type"`
	Items     []Item           `json:"items"`
	CreatedAt time.Time        `json:"createdAt"`
	UpdatedAt time.Time        `json:"updateAt"`
}

func ToList(l dbmodel.List) List {
	ll := List{
		ID:        l.ID,
		Title:     l.Title,
		Type:      l.Type,
		CreatedAt: l.CreatedAt,
		UpdatedAt: l.UpdatedAt,
		Items:     []Item{},
	}

	for _, i := range l.Items {
		ll.Items = append(ll.Items, ToItem(i))
	}

	return ll
}

func ToAllLists(l []dbmodel.List) []List {
	ll := []List{}

	for _, dbl := range l {
		ll = append(ll, ToList(dbl))
	}
	return ll
}

type CreateList struct {
	Title string
	Type  dbmodel.ListType
}

type UpdateList struct {
	Title string
	Type  dbmodel.ListType
}
