package dbmodel

import (
	"gorm.io/gorm"
)

type Item struct {
	gorm.Model
	ListID    string
	Text      string
	State     string
	SortIndex int
}
