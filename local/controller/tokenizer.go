package controller

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/Barokah-AI/BackEnd/model"
)

// Baca file vocab.txt dan return map dari kata ke index
func readVocab(filePath string) (map[string]int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open vocab file: %v", err)
	}
	defer file.Close()

	vocab := make(map[string]int)
	scanner := bufio.NewScanner(file)
	index := 0
	for scanner.Scan() {
		vocab[scanner.Text()] = index
		index++
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading vocab file: %v", err)
	}

	return vocab, nil
}

// Baca file tokenizer_config.json dan return struct konfigurasi
func readTokenizerConfig(filePath string) (*model.TokenizerConfig, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open tokenizer config file: %v", err)
	}
	defer file.Close()

	var config model.TokenizerConfig
	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("failed to decode tokenizer config: %v", err)
	}

	return &config, nil
}

// Tokenisasi text berdasarkan vocab dan konfigurasi tokenizer
func tokenize(text string, vocab map[string]int, config *model.TokenizerConfig) ([]int, error) {
	if config.DoLowerCase {
		text = strings.ToLower(text)
	}

	// Split the text into words
	words := strings.Fields(text)

	var tokens []int
	// Tambahkan CLS token di awal
	if token, ok := vocab[config.ClsToken]; ok {
		tokens = append(tokens, token)
	} else {
		return nil, fmt.Errorf("CLS token not found in vocab")
	}

	for _, word := range words {
		if token, ok := vocab[word]; ok {
			tokens = append(tokens, token)
		} else {
			if unkToken, ok := vocab[config.UnkToken]; ok {
				tokens = append(tokens, unkToken)
			} else {
				return nil, fmt.Errorf("word not found in vocab: %s and UNK token not found", word)
			}
		}
	}

	// Tambahkan SEP token di akhir
	if token, ok := vocab[config.SepToken]; ok {
		tokens = append(tokens, token)
	} else {
		return nil, fmt.Errorf("SEP token not found in vocab")
	}

	return tokens, nil
}
