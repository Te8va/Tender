package errwriter

import (
	"encoding/json"
	"net/http"

	"git.codenrock.com/cnrprod1725727333-user-88349/zadanie-6105/internal/tender/domain"
)

func RespondWithError(w http.ResponseWriter, statusCode int, errMessage string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(domain.JSONError{Err: errMessage})
}
