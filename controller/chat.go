package controller

// import (
// 	"encoding/json"
// 	"log"
// 	"net/http"
// 	"strconv"

// 	"github.com/Barokah-AI/BackEnd/config"
// 	"github.com/Barokah-AI/BackEnd/helper"
// 	"github.com/Barokah-AI/BackEnd/model"
// 	"github.com/go-resty/resty/v2"
// )

func Chat(respw http.ResponseWriter, req *http.Request, tokenmodel string) {
    var chat model.AIRequest