package helper

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	"strconv"
)

// LoadDataset loads the dataset from GCS bucket and returns a map of label to question-answer pair
func LoadDatasetGCS(bucketName, objectName string) (map[string][]string, error) {
    data, err := ReadFileFromGCS(bucketName, objectName)
    if err != nil {
        return nil, fmt.Errorf("failed to read dataset file: %v", err)
    }

    r := bytes.NewReader(data)
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