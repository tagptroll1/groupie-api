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
	api.Use(middleware.SetHeader("Content-Type", "application/json"))

	api.Route("/lists", func(r chi.Router) {
		lists := service.NewListService(db)
		items := service.NewItemService(db)

		r.Get("/", lists.GetAllLists)
		r.Post("/", lists.Create)

		r.Route("/{listkey}", func(r chi.Router) {
			r.Get("/", lists.Get)
			r.Put("/", lists.Update)
			r.Delete("/", lists.Delete)

			r.Route("/items", func(r chi.Router) {
				r.Get("/", items.ListItems)
				r.Post("/", items.Create)

				r.Route("/{item}", func(r chi.Router) {
					r.Get("/", items.Get)
					r.Put("/", items.Update)
					r.Delete("/", items.Delete)

				})

			})
		})
	})

	r.Mount("/api", api)
	return r
}
