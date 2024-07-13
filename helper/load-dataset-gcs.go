package helper

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"strconv"

	"cloud.google.com/go/storage"
)

// LoadDataset loads the dataset from GCS bucket and returns a map of label to question-answer pair
func LoadDatasetGCS(bucketName, objectName string) (map[string][]string, error) {
    ctx := context.Background()
    client, err := storage.NewClient(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to create storage client: %v", err)
    }
    defer client.Close()

    bucket := client.Bucket(bucketName)
    object := bucket.Object(objectName)
    r, err := object.NewReader(ctx)
    if err != nil {
        return nil, fmt.Errorf("failed to create reader for object: %v", err)
    }
    defer r.Close()

    reader := csv.NewReader(r)
    reader.Comma = '|' // Set the delimiter to pipe
    records, err := reader.ReadAll()
    if err != nil {
        return nil, fmt.Errorf("failed to read dataset file: %v", err)
    }

    labelToQA := make(map[string][]string)
    for i, record := range records {
        if len(record) != 2 {
            log.Printf("Skipping invalid record at line %d: %v\n", i+1, record)
            continue
        }
        label := "LABEL_" + strconv.Itoa(i)
        labelToQA[label] = record
    }
    return labelToQA, nil
}