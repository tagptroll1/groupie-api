package service

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tagptroll1/groupie-api/model/dbmodel"
	"gorm.io/gorm"
)

type ListService struct {
	db *gorm.DB
}

func NewListService(db *gorm.DB) *ListService {
	return &ListService{
		db: db,
	}
}

func (s *ListService) GetAllLists(w http.ResponseWriter, r *http.Request) {
	var lists []dbmodel.List
	s.db.Find(&lists)
	json.NewEncoder(w).Encode(lists)
}

func (s *ListService) GetList(w http.ResponseWriter, r *http.Request) {
	listId := chi.URLParam(r, "listkey")
	var list dbmodel.List
	err := s.db.Model(&dbmodel.List{}).
		Preload("Items").
		Find(&list, listId).
		Error

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(list)
}

func (s *ListService) GetItem(w http.ResponseWriter, r *http.Request) {
	itemID := chi.URLParam(r, "item")
	listId := chi.URLParam(r, "listkey")

	var item dbmodel.Item
	s.db.Where("list_id = ?", listId).Find(&item, itemID)
	json.NewEncoder(w).Encode(item)
}
