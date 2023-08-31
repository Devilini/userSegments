package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Handler struct {
}

func (h *Handler) InitRoutes(router *httprouter.Router) {
	//router.HandlerFunc(http.MethodGet, "/api/user", GetUser())
}

func responseJson(w http.ResponseWriter, data any) {
	jsonResp, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
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
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	w.Write(jsonResp)
}
