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
		return nil, fmt.Errorf("gagal dalam membaca file dataset: %v", err)
	}

	label_to_qa := make(map[string][]string)
	for i, record := range records {
		if len(record) != 2 {
			log.Printf("Melewati record yang tidak valid pada baris %d: %v\n", i+1, record)
			continue
		}
		label := "LABEL_" + strconv.Itoa(i)
		label_to_qa[label] = record
	}
	return label_to_qa, nil
}
