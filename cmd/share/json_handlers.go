package share

import (
	"encoding/json"
	"net/http"
)

func JsonResponse(response interface{}, w http.ResponseWriter) error {
	json, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)

	return nil
}
