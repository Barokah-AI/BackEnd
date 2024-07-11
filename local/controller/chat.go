package controller

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/Barokah-AI/BackEnd/config"
	"github.com/Barokah-AI/BackEnd/helper"
	"github.com/Barokah-AI/BackEnd/model"
	"github.com/go-resty/resty/v2"
)

// LoadDataset loads the dataset from the given CSV file and returns a map of label to question-answer pair
func LoadDataset(filePath string) (map[string][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open dataset file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
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

// NormalizeText normalizes the text by converting to lowercase and removing punctuation
func NormalizeText(text string) string {
	// Convert to lowercase
	text = strings.ToLower(text)
	// Remove punctuation
	text = strings.Map(func(r rune) rune {
		if strings.ContainsRune("abcdefghijklmnopqrstuvwxyz0123456789 ", r) {
			return r
		}
		return -1
	}, text)
	return text
}

