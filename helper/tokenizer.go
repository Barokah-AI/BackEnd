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

// Fungsi untuk membaca tokenizer config dari GCS
func ReadTokenizerConfigFromGCS(bucketName, fileName string) (map[string]interface{}, error) {
	data, err := ReadFileFromGCS(bucketName, fileName)
	if err != nil {
		return nil, err
	}

	var config map[string]interface{}
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling tokenizer config: %v", err)
	}

	return config, nil
}

// Baca file vocab.txt dan return map dari kata ke index
func ReadVocab(filePath string) (map[string]int, error) {
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
func ReadTokenizerConfig(filePath string) (*model.TokenizerConfig, error) {
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
func Tokenize(text string, vocab map[string]int, config *model.TokenizerConfig) ([]int, error) {
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

func Tokenize2(text string, vocab map[string]int, tokenizerConfig map[string]interface{}) ([]int, error) {
	// Simpan token ke ID dalam map
	tokenToID := make(map[string]int)
	for k, v := range vocab {
		id, err := strconv.Atoi(strconv.Itoa(v))
		if err != nil {
			return nil, fmt.Errorf("invalid token ID in vocab: %v", err)
		}
		tokenToID[k] = id
	}

	// Tokenisasi sederhana berdasarkan spasi
	words := strings.Fields(text)
	tokens := []int{}

	for _, word := range words {
		id, exists := tokenToID[word]
		if !exists {
			id = tokenToID[tokenizerConfig["unk_token"].(string)]
		}
		tokens = append(tokens, id)
	}

	// Tambahkan token CLS dan SEP
	clsToken := tokenToID[tokenizerConfig["cls_token"].(string)]
	sepToken := tokenToID[tokenizerConfig["sep_token"].(string)]
	tokens = append([]int{clsToken}, tokens...)
	tokens = append(tokens, sepToken)

	return tokens, nil
}
