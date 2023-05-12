package service

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tagptroll1/groupie-api/model/dbmodel"
	"github.com/tagptroll1/groupie-api/model/rest"
	"gorm.io/gorm"
)

type ItemService struct {
	db *gorm.DB
}

func NewItemService(db *gorm.DB) *ItemService {
	return &ItemService{db: db}
}

func (s *ItemService) ListItems(w http.ResponseWriter, r *http.Request) {
	listId := chi.URLParam(r, "listkey")

	var items []dbmodel.Item
	err := s.db.WithContext(r.Context()).Find(&items, "list_id = ?", listId).Error

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(items)
}

func (s *ItemService) Get(w http.ResponseWriter, r *http.Request) {
	itemID := chi.URLParam(r, "item")
	listId := chi.URLParam(r, "listkey")

	var item dbmodel.Item
	err := s.db.
		WithContext(r.Context()).
		Where("list_id = ?", listId).
		Find(&item, "id", itemID).
		Error

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(item)
}

func (s *ItemService) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	itemID := chi.URLParam(r, "item")
	listId := chi.URLParam(r, "listkey")

	var item *dbmodel.Item
	err := s.db.
		WithContext(ctx).
		Where("list_id = ?", listId).
		Find(&item, "id", itemID).
		Error

	if err != nil || item == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.db.WithContext(ctx).Where("id = ?", itemID).Delete(&dbmodel.Item{}).Error

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *ItemService) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	listId := chi.URLParam(r, "listkey")

	var list *dbmodel.List
	err := s.db.WithContext(ctx).
		Where("id = ?", listId).
		First(&list).Error

	if err != nil || list == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var item *rest.CreateItem
	err = json.NewDecoder(r.Body).Decode(&item)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	dbitem := &dbmodel.Item{
		Text:      item.Text,
		ListID:    listId,
		State:     item.State,
		SortIndex: item.SortIndex,
	}
	err = s.db.Model(&dbitem).
		WithContext(ctx).
		Save(&dbitem).Error

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}

func (s *ItemService) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	itemID := chi.URLParam(r, "item")
	listId := chi.URLParam(r, "listkey")

	var list *dbmodel.List
	err := s.db.WithContext(ctx).
		Where("id = ?", listId).
		First(&list).Error

	if err != nil || list == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var item map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&item)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	delete(item, "updated_at")

	err = s.db.
		WithContext(ctx).
		Model(&dbmodel.Item{}).
		Where("id = ?", itemID).
		Omit("list_id", "id", "created_at").
		Updates(&item).Error

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(item)
}
