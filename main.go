package barokahai

import (
	routes "github.com/Barokah-AI/BackEnd/url"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

// Fungsi init akan dipanggil ketika package ini diinisialisasi
func init() {
	functions.HTTP("WebHook", routes.URL)
}