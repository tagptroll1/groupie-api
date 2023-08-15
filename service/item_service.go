package service

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/tagptroll1/groupie-api/model/dbmodel"
	"github.com/tagptroll1/groupie-api/model/rest"

	"github.com/go-chi/chi/v5"
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

	var list dbmodel.List
	res := s.db.WithContext(r.Context()).
		Model(&dbmodel.List{}).
		Preload("Items").
		Where("id = ?", listId).
		Find(&list)

	if res.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if res.RowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(rest.ToList(list).Items)
}

func (s *ItemService) Get(w http.ResponseWriter, r *http.Request) {
	itemID := chi.URLParam(r, "item")
	listId := chi.URLParam(r, "listkey")

	var item dbmodel.Item
	res := s.db.
		WithContext(r.Context()).
		Where("list_id = ?", listId).
		Find(&item, "id", itemID)

	if res.Error != nil {
		log.Print(res.Error)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if res.RowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(rest.ToItem(item))
}

func (s *ItemService) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	itemID := chi.URLParam(r, "item")
	listId := chi.URLParam(r, "listkey")

	var item dbmodel.Item
	res := s.db.WithContext(ctx).
		Limit(1).
		Where("list_id = ?", listId).
		Find(&item, "id", itemID)

	if res.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if res.RowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	res = s.db.WithContext(ctx).Where("id = ?", itemID).Delete(&dbmodel.Item{})

	if res.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if res.RowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *ItemService) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	listId := chi.URLParam(r, "listkey")

	var list *dbmodel.List
	res := s.db.WithContext(ctx).
		Limit(1).
		Find(&list, "id = ?", listId)

	if res.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if res.RowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var item rest.CreateItem
	err := json.NewDecoder(r.Body).Decode(&item)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dbitem := dbmodel.Item{
		Text:      item.Text,
		ListID:    listId,
		State:     item.State,
		SortIndex: item.SortIndex,
	}
	res = s.db.Model(&dbitem).
		WithContext(ctx).
		Save(&dbitem)

	if res.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(rest.ToItem(dbitem))
}

func (s *ItemService) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	itemID := chi.URLParam(r, "item")
	listId := chi.URLParam(r, "listkey")

	var list *dbmodel.List
	res := s.db.WithContext(ctx).
		Limit(1).
		Find(&list, "id = ?", listId)

	if res.Error != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if res.RowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var item map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&item)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	delete(item, "updated_at")

	res = s.db.
		WithContext(ctx).
		Model(&dbmodel.Item{}).
		Where("id = ?", itemID).
		Omit("list_id", "id", "created_at").
		Updates(&item)

	if res.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if res.RowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(item)
}
