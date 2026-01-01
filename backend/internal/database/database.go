package database

import (
	"log"
	"os"

	"github.com/supabase-community/postgrest-go"
	storage_go "github.com/supabase-community/storage-go"
)

// AppClients holds the PostgREST and Storage clients
type AppClients struct {
	Postgrest *postgrest.Client
	Storage   *storage_go.Client
}

// InitDB initializes and returns a custom struct containing PostgREST and Storage clients
func InitDB() (*AppClients, error) {
	supabaseRestURL := os.Getenv("SUPABASE_REST_URL")
	supabaseStorageURL := os.Getenv("SUPABASE_STORAGE_URL")
	supabaseKey := os.Getenv("SUPABASE_KEY")

	log.Printf("DEBUG: Supabase REST URL: %s", supabaseRestURL)
	log.Printf("DEBUG: Supabase Storage URL: %s", supabaseStorageURL)
	log.Printf("DEBUG: Supabase Key (first 10 chars): %s", supabaseKey[:min(10, len(supabaseKey))])

	if supabaseRestURL == "" || supabaseStorageURL == "" || supabaseKey == "" {
		log.Fatal("SUPABASE_REST_URL, SUPABASE_STORAGE_URL, or SUPABASE_KEY is not set in .env or is empty")
	}

	// Initialize PostgREST client
	postgrestClient := postgrest.NewClient(supabaseRestURL, "", map[string]string{
		"apikey":        supabaseKey,
		"Authorization": "Bearer " + supabaseKey,
	})

	if postgrestClient.ClientError != nil {
		return nil, postgrestClient.ClientError
	}

	// Initialize Storage client
	log.Printf("DEBUG: Initializing Storage client with URL: %s, Key prefix: %s", supabaseStorageURL, supabaseKey[:min(10, len(supabaseKey))])
	storageClient := storage_go.NewClient(supabaseStorageURL, supabaseKey, map[string]string{
		"Authorization": "Bearer " + supabaseKey,
	})

	// You can optionally test connections here if needed
	// For example, try to list buckets with storageClient
	_, err := storageClient.GetBucket("balkanid-file-storage")
	if err != nil {
		log.Printf("DEBUG: Bucket 'balkanid-file-storage' not found, attempting to create it...")
		_, err = storageClient.CreateBucket("balkanid-file-storage", storage_go.BucketOptions{
			Public: false,
		})
		if err != nil {
			log.Printf("FATAL: Could not create bucket 'balkanid-file-storage'")
			return nil, err
		}
		log.Printf("DEBUG: Bucket 'balkanid-file-storage' created successfully.")
	} else {
		log.Printf("DEBUG: Bucket 'balkanid-file-storage' already exists.")
	}

	return &AppClients{
		Postgrest: postgrestClient,
		Storage:   storageClient,
	}, nil
}
