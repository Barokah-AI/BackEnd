package helper

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
)

// LoadDataset loads the dataset from the given CSV file and returns a map of label to question-answer pair
func LoadDatasetLocal(file_path string) (map[string][]string, error) {
	file, err := os.Open(file_path)
	if err != nil {
		return nil, fmt.Errorf("gagal dalam membuka file dataset: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '|' // Set the delimiter to pipe
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("gagal dalam membaca dataset: %v", err)
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
