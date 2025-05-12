package main

import (
	"bookmarker/internal/dbutil"
	"bookmarker/internal/repositories"
	"bookmarker/internal/services"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// importPinboard runs the import-pinboard command
func importPinboard(filename string) {
	db, err := dbutil.OpenPostgresDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	bookmarkRepo := repositories.NewBookmarkRepository(db)
	tagRepo := repositories.NewTagRepository(db)
	bookmarkService := services.NewBookmarkServiceWithTags(bookmarkRepo, tagRepo)
	importService := services.NewPinboardImportService(bookmarkService)

	filePath := filepath.Join("../../data/import", filename)
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Failed to open import file: %v", err)
	}
	defer f.Close()

	err = importService.ImportFromJSON(f)
	if err != nil {
		log.Fatalf("Import failed: %v", err)
	}
	fmt.Println("Pinboard import completed successfully.")
}
