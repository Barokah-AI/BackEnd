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