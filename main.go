package main

import (
	"context"
	"net/http"
	"os"

	"github.com/tagptroll1/groupie-api/model/dbmodel"
	"github.com/tagptroll1/groupie-api/router"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	ctx := context.Background()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3333"
	}

	cs := os.Getenv("DATABASE_URL")
	db, err := setupDatabase(cs)

	if err != nil {
		panic(err)
	}

	r := router.New(ctx, db)

	http.ListenAndServe(":"+port, r)
}

func setupDatabase(cs string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cs), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&dbmodel.List{})
	db.AutoMigrate(&dbmodel.Item{})

	return db, nil
}
