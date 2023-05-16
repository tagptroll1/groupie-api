package router

import (
	"context"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/tagptroll1/groupie-api/service"
	"gorm.io/gorm"
)

func New(ctx context.Context, db *gorm.DB) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	api := chi.NewRouter()
	api.Use(middleware.SetHeader("Content-Type", "application/json"))

	if db == nil {
		return r
	}

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
