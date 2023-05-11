package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3333"
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(60 * time.Second))

	api := chi.NewRouter()
	api.Route("/lists", func(r chi.Router) {
		r.Get("/", getAllLists)
		r.Route("/{list}", func(r chi.Router) {
			r.Get("/", getList)
			r.Route("/{item}", func(r chi.Router) {
				r.Get("/", getItem)
			})
		})

	})

	r.Mount("/api", api)

	http.ListenAndServe(":"+port, r)
}

func getAllLists(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`{"lists": [{name: "yehaw"}]}`))
}

func getList(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf(`{"name": "%s","list": []}`, chi.URLParam(r, "list"))))
}

func getItem(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(fmt.Sprintf(`{"name": "%s","watched": true}`, chi.URLParam(r, "item"))))
}
