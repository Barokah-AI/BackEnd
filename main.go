package main

import (
	routes "github.com/Barokah-AI/BackEnd/url"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	functions.HTTP("WebHook", routes.URL)
}