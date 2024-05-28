package main

import (
	"context"
	"log"
	"woonbeaj/segments/config"
	"woonbeaj/segments/internal/storage"
)

func main() {
	cfg := config.MustLoad(`config\config.yaml`)
	db, err := storage.New(cfg.StorageUrl)
	if err != nil {
		log.Fatalf("connect db failed: %v", err)
	}
	defer db.Close(context.Background())
	db.Clean()
}