package controller

import (
	"encoding/json"
	"net/http"
)

func responseJson(w http.ResponseWriter, data any) {
	resp := make(map[string]any)
	resp["data"] = data
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}

func errorResponseJson(w http.ResponseWriter, text string) {
	resp := make(map[string]string)
	resp["error"] = text
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	w.Write(jsonResp)
}
