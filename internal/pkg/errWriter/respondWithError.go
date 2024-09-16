package errwriter

import (
	"encoding/json"
	"net/http"

	"github.com/Te8va/Tender/internal/tender/domain"
)

func RespondWithError(w http.ResponseWriter, statusCode int, errMessage string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(domain.JSONError{Err: errMessage})
}
