package helper

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Barokah-AI/BackEnd/model"
)

// Fungsi untuk membaca vocab dari GCS
func ReadVocabFromGCS(bucketName, fileName string) (map[string]int, error) {
	data, err := ReadFileFromGCS(bucketName, fileName)
	if err != nil {
		return nil, err
	}

	vocab := make(map[string]int)
	scanner := bufio.NewScanner(bytes.NewReader(data))
	index := 0
	for scanner.Scan() {
		line := scanner.Text()
		token := strings.TrimSpace(line)
		if token == "" {
			return nil, fmt.Errorf("invalid line in vocab file: %s", line)
		}
		vocab[token] = index
		index++
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading vocab file: %v", err)
	}

	return vocab, nil
}
