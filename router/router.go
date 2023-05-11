package router

import (
	"context"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/tagptroll1/groupie-api/service"
	"gorm.io/gorm"
)

func New(ctx context.Context, db *gorm.DB) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(60 * time.Second))

	api := chi.NewRouter()

	api.Route("/lists", func(r chi.Router) {
		lists := service.NewListService(db)

		r.Get("/", lists.GetAllLists)
		r.Route("/{listkey}", func(r chi.Router) {
			r.Get("/", lists.GetList)
			r.Route("/{item}", func(r chi.Router) {
				r.Get("/", lists.GetItem)
			})
		})
	})

	r.Mount("/api", api)
	return r
}
