package dbmodel

import (
	"time"

	"gorm.io/gorm"
)

type ListType string

const (
	Thumbs ListType = "thumbs"
)

type List struct {
	ID        string `gorm:"primarykey"`
	Title     string
	Type      ListType
	Items     []Item
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
