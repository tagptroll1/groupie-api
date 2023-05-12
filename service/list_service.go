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
	s.db.Find(&lists)
	json.NewEncoder(w).Encode(lists)
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
	err = s.db.Model(&dbmodel.List{}).
		Create(&dbList).Error

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dbList)
}

func (s *ListService) Get(w http.ResponseWriter, r *http.Request) {
	listId := chi.URLParam(r, "listkey")
	var list dbmodel.List
	err := s.db.Model(&dbmodel.List{}).
		Preload("Items").
		Find(&list, "id", listId).
		Error

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(list)
}

func (s *ListService) Update(w http.ResponseWriter, r *http.Request) {
	listId := chi.URLParam(r, "listkey")

	var list map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&list)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.db.Model(&dbmodel.List{ID: listId}).Select("title").Updates(list).Error

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}
