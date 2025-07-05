package clients

import (
	"os"
	"path/filepath"

	storage_go "github.com/supabase-community/storage-go"
)

// SupabaseDatastoreClient uploads files to Supabase Storage
// Requires SUPABASE_URL, SUPABASE_SERVICE_KEY, SUPABASE_BUCKET env vars
// Uses the Supabase Storage REST API

type SupabaseDatastoreClient struct {
	Url        string
	ServiceKey string
	Bucket     string
}

func NewSupabaseDatastoreClient() *SupabaseDatastoreClient {
	return &SupabaseDatastoreClient{
		Url:        os.Getenv("SUPABASE_S3_URL"),
		ServiceKey: os.Getenv("SUPABASE_SERVICE_KEY"),
		Bucket:     os.Getenv("SUPABASE_BUCKET"),
	}
}

// UploadFile uploads a file at filePath to Supabase Storage using the storage-go client
func (c *SupabaseDatastoreClient) UploadFile(filePath string) error {
	client := storage_go.NewClient(c.Url, c.ServiceKey, nil)
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	fileName := filepath.Base(filePath)
	// Upload file using Put method
	_, err = client.UploadFile(c.Bucket, fileName, file)
	if err != nil {
		return err
	}
	return nil
}