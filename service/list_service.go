package service

import (
	"encoding/json"
	"net/http"

	"github.com/tagptroll1/groupie-api/model/dbmodel"
	"github.com/tagptroll1/groupie-api/model/rest"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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

// TODO: Don't expose this later
func (s *ListService) GetAllLists(w http.ResponseWriter, r *http.Request) {
	var lists []dbmodel.List
	s.db.WithContext(r.Context()).Find(&lists)

	if len(lists) == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	json.NewEncoder(w).Encode(rest.ToAllLists(lists))
}

func (s *ListService) Create(w http.ResponseWriter, r *http.Request) {
	var list rest.CreateList
	err := json.NewDecoder(r.Body).Decode(&list)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dbList := dbmodel.List{
		ID:    uuid.New().String(),
		Title: list.Title,
		Type:  list.Type,
	}
	err = s.db.WithContext(r.Context()).
		Model(&dbmodel.List{}).
		Create(&dbList).Error

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(rest.ToList(dbList))
}

func (s *ListService) Get(w http.ResponseWriter, r *http.Request) {
	listId := chi.URLParam(r, "listkey")

	var list dbmodel.List
	res := s.db.Model(&dbmodel.List{}).
		WithContext(r.Context()).
		Preload("Items").
		Find(&list, "id", listId)

	if res.Error != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if res.RowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(rest.ToList(list))
}

func (s *ListService) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	listId := chi.URLParam(r, "listkey")

	var list *dbmodel.List
	res := s.db.Model(&dbmodel.List{}).
		WithContext(ctx).
		Preload("Items").
		Find(&list, "id", listId)

	if res.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if res.RowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	res = s.db.WithContext(ctx).Where("list_id = ?", listId).Delete(&dbmodel.Item{})

	if res.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to delete items"))
		return
	}

	res = s.db.WithContext(ctx).Delete(&dbmodel.List{ID: listId})

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

func (s *ListService) Update(w http.ResponseWriter, r *http.Request) {
	listId := chi.URLParam(r, "listkey")

	var list map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&list)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res := s.db.WithContext(r.Context()).
		Model(&dbmodel.List{ID: listId}).
		Select("title").
		Updates(list)

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
