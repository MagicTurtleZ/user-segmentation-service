package main

import (
	"context"
	"log"
	"woonbeaj/segments/config"
	"woonbeaj/segments/internal/handlers/audit"
	"woonbeaj/segments/internal/handlers/segments/create"
	segdelete "woonbeaj/segments/internal/handlers/segments/delete"
	addsegments "woonbeaj/segments/internal/handlers/user/add_segments"
	showsegments "woonbeaj/segments/internal/handlers/user/show_segments"
	"woonbeaj/segments/internal/storage"

	"github.com/gofiber/fiber/v2"
)

func main() {
	cfg := config.MustLoad(`config\config.yaml`)
	db, err := storage.New(cfg.StorageUrl)
	if err != nil {
		log.Fatalf("connect db failed: %v", err)
	}
	defer db.Close(context.Background())

	app := fiber.New()
	app.Post("/create-seg", create.NewSegment(db))
	app.Post("/delete-seg", segdelete.RemoveSegment(db))
	app.Post("/add-user-segments", addsegments.AddSegmentsToUser(db))
	app.Post("/show-user-segments", showsegments.AddSegmentsToUser(db))
	app.Post("/report", audit.ReportCSV(db))
	err = app.Listen(cfg.Address)
	if err != nil {
		log.Fatal(err)
	}
}
