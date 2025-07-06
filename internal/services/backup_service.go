package services

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"bookmarker/internal/clients"
)

// BackupPostgresDB runs pg_dump using environment variables and saves output to a timestamped file.
func BackupPostgresDB() error {
	pgUser := os.Getenv("DB_USER")
	pgPassword := os.Getenv("DB_PASS")
	pgDB := os.Getenv("DB_NAME")
	pgHost := os.Getenv("DB_HOST")
	pgPort := os.Getenv("DB_PORT")
	if pgUser == "" || pgPassword == "" || pgDB == "" {
		return fmt.Errorf("DB_USER, DB_PASS, and DB_NAME must be set in environment")
	}
	if pgHost == "" {
		pgHost = "localhost"
	}
	if pgPort == "" {
		pgPort = "5432"
	}
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("data/backup/backup_%s.sql", timestamp)
	cmd := exec.Command(
		"pg_dump",
		"-h", pgHost,
		"-p", pgPort,
		"-U", pgUser,
		pgDB,
	)
	cmd.Env = append(os.Environ(), "PGPASSWORD="+pgPassword)
	outFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer outFile.Close()
	cmd.Stdout = outFile
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	// Compress the SQL file before upload
	gzFilename := filename + ".gz"
	gzCmd := exec.Command("gzip", "-c", filename)
	gzOutFile, err := os.Create(gzFilename)
	if err != nil {
		return err
	}
	defer gzOutFile.Close()
	gzCmd.Stdout = gzOutFile
	gzCmd.Stderr = os.Stderr
	if err := gzCmd.Run(); err != nil {
		return err
	}
	// Optionally remove the uncompressed file to save space
	os.Remove(filename)

	// Upload to Supabase after backup (use compressed file)
	supabaseClient := clients.NewSupabaseDatastoreClient()
	uploadErr := supabaseClient.UploadFile(gzFilename)
	if uploadErr != nil {
		return fmt.Errorf("backup created but upload to supabase failed: %w", uploadErr)
	}

	return nil
}