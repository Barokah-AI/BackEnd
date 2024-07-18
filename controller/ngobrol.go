package controller

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"strconv"
// 	"strings"

// 	"github.com/Barokah-AI/BackEnd/config"
// 	"github.com/Barokah-AI/BackEnd/helper"
// 	"github.com/Barokah-AI/BackEnd/model"
// )

func Ngobrol(respw http.ResponseWriter, req *http.Request, tokenmodel string) {
    var chat model.AIRequest

    err := json.NewDecoder(req.Body).Decode(&chat)
    if err != nil {
        helper.ErrorResponse(respw, req, http.StatusBadRequest, "Permintaan Tidak Valid", "error saat membaca isi permintaan: "+err.Error())
        return
    }

	if chat.Prompt == "" {
        helper.ErrorResponse(respw, req, http.StatusBadRequest, "Permintaan Tidak Valid", "masukin pertanyaan dulu ya kakak 🤗")
        return
    }

	bucketName := config.GetEnv("GCS_BUCKET_NAME")
    vocabObjectName := config.GetEnv("GCS_VOCAB_FILE")
    tokenizerConfigName := config.GetEnv("GCS_TOKENIZER_CONFIG_FILE")

	 // Read and use the tokenizer
	 vocab, err := helper.ReadVocabFromGCS(bucketName, vocabObjectName)
	 if err != nil {
		 helper.ErrorResponse(respw, req, http.StatusInternalServerError, "Kesalahan Server Internal", "tidak bisa membaca vocab: "+err.Error())
		 return
	 }
 
	 tokenizerConfig, err := helper.ReadTokenizerConfigFromGCS(bucketName, tokenizerConfigName)
	 if err != nil {
		 helper.ErrorResponse(respw, req, http.StatusInternalServerError, "Kesalahan Server Internal", "tidak bisa membaca konfigurasi tokenizer: "+err.Error())
		 return
	 }

	 tokens, err := helper.Tokenize2(chat.Prompt, vocab, tokenizerConfig)
	if err != nil {
		helper.ErrorResponse(respw, req, http.StatusInternalServerError, "Kesalahan Server Internal", "error saat melakukan tokenisasi: "+err.Error())
		return
	}

	// Convert tokens to string for API call
	tokensStr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(tokens)), " "), "[]")

	// Call Hugging Face API with tokenized prompt
	label, score, err := helper.CallHuggingFaceAPI(tokensStr)
	if err != nil {
		helper.ErrorResponse(respw, req, http.StatusInternalServerError, "Kesalahan Server Internal", "model sedang diload: "+err.Error())
		return
	}