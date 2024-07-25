package helper

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/storage"
)

// Fungsi untuk membaca file dari GCS
func ReadFileFromGCS(bucketName, fileName string) ([]byte, error) {
	ctx := context.Background()

	clients, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %v", err)
	}
	defer clients.Close()

	buckets := clients.Bucket(bucketName)
	obj := buckets.Object(fileName)
	r, err := obj.NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create reader: %v", err)
	}
	defer r.Close()

	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read data: %v", err)
	}

	return data, nil
}
