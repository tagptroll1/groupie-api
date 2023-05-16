package main

import (
	"context"
	"embed"
	_ "embed"
	"log"
	"net/http"
	"os"

	"github.com/tagptroll1/groupie-api/model/dbmodel"
	"github.com/tagptroll1/groupie-api/router"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/prometheus"
)

//go:embed swagger/*
var swagger embed.FS

func main() {
	ctx := context.Background()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3333"
	}

	cs := os.Getenv("DATABASE_URL")

	db, err := setupDatabase(cs)

	if err != nil {
		log.Println(err)
	}

	r := router.New(ctx, db)

	fs := http.FileServer(http.FS(swagger))
	r.Handle("/swagger", http.RedirectHandler("/swagger/", http.StatusPermanentRedirect))
	r.Handle("/swagger/*", http.StripPrefix("/swagger/", fs))

	http.ListenAndServe(":"+port, r)
}

func setupDatabase(cs string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cs), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	db.Use(prometheus.New(prometheus.Config{
		DBName:          "groupie",
		RefreshInterval: 15,
		StartServer:     true,
		HTTPServerPort:  9091,
		MetricsCollector: []prometheus.MetricsCollector{
			&prometheus.Postgres{
				VariableNames: []string{"Threads_running"},
			},
		},
	}))

	db.AutoMigrate(&dbmodel.List{})
	db.AutoMigrate(&dbmodel.Item{})

	return db, nil
}
