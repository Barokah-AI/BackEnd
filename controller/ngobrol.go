package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Barokah-AI/BackEnd/config"
	"github.com/Barokah-AI/BackEnd/helper"
	"github.com/Barokah-AI/BackEnd/model"
)

func Ngobrol(respwrt http.ResponseWriter, req *http.Request, tokenmodel string) {
	var chat model.AIRequest

	err := json.NewDecoder(req.Body).Decode(&chat)
	if err != nil {
		helper.ErrorResponse(respwrt, req, http.StatusBadRequest, "Permintaan Tidak Valid", "error saat membaca isi permintaan: "+err.Error())
		return
	}

	if chat.Prompt == "" {
		helper.ErrorResponse(respwrt, req, http.StatusBadRequest, "Permintaan Tidak Valid", "masukin pertanyaan dulu ya kakak 🤗")
		return
	}

	bucket_name := config.GetEnv("GCS_BUCKET_NAME")
	vocab_object_name := config.GetEnv("GCS_VOCAB_FILE")
	tokenizer_config_name := config.GetEnv("GCS_TOKENIZER_CONFIG_FILE")

	// Read and use the tokenizer
	vocab, err := helper.ReadVocabFromGCS(bucket_name, vocab_object_name)
	if err != nil {
		helper.ErrorResponse(respwrt, req, http.StatusInternalServerError, "Kesalahan Server Internal", "tidak bisa membaca vocab: "+err.Error())
		return
	}

	tokenizerConfig, err := helper.ReadTokenizerConfigFromGCS(bucket_name, tokenizer_config_name)
	if err != nil {
		helper.ErrorResponse(respwrt, req, http.StatusInternalServerError, "Kesalahan Server Internal", "tidak bisa membaca konfigurasi tokenizer: "+err.Error())
		return
	}

	tokens, err := helper.Tokenize2(chat.Prompt, vocab, tokenizerConfig)
	if err != nil {
		helper.ErrorResponse(respwrt, req, http.StatusInternalServerError, "Kesalahan Server Internal", "error saat melakukan tokenisasi: "+err.Error())
		return
	}

	// Convert tokens to string for API call
	tokens_string := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(tokens)), " "), "[]")

	// Call Hugging Face API with tokenized prompt
	label, score, err := helper.CallHuggingFaceAPI(tokens_string)
	if err != nil {
		helper.ErrorResponse(respwrt, req, http.StatusInternalServerError, "Kesalahan Server Internal", "model sedang diload: "+err.Error())
		return
	}

	// Load the dataset from GCS
	bucket_name = config.GetEnv("GCS_BUCKET_NAME")
	object_name := config.GetEnv("GCS_OBJECT_NAME")

	label_t0_qa, err := helper.LoadDatasetGCS(bucket_name, object_name)
	if err != nil {
		helper.ErrorResponse(respwrt, req, http.StatusInternalServerError, "Kesalahan Server Internal", "kesalahan server: tidak bisa memuat dataset: "+err.Error())
		return
	}

	// Get the answer corresponding to the best label
	record, ok := label_t0_qa[label]
	if !ok {
		helper.ErrorResponse(respwrt, req, http.StatusInternalServerError, "Kesalahan Server Internal", "kesalahan server: label tidak ditemukan dalam dataset")
		return
	}

	answer := record[1]

	helper.WriteJSON(respwrt, http.StatusOK, map[string]string{
		"prompt":   chat.Prompt,
		"response": answer,
		"label":    label,
		"score":    strconv.FormatFloat(score, 'f', -1, 64),
	})
}
