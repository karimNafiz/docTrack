package pdf

import (
	"net/http"
)

func CreatePDFHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	// var req struct {
	// }
}
