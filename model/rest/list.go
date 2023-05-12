package rest

import "github.com/tagptroll1/groupie-api/model/dbmodel"

type CreateList struct {
	Title string
	Type  dbmodel.ListType
}

type UpdateList struct {
	ID    string
	Title string
	Type  dbmodel.ListType
}
