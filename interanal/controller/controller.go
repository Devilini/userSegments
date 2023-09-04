package controller

import (
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"userSegments/interanal/apperror"
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

func errorResponseJson(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	var appErr *apperror.AppError
	if err != nil {
		logrus.Error(err.Error())
		if errors.As(err, &appErr) {
			log.Println(errors.Is(err, apperror.ErrNotFound))

			if errors.Is(err, apperror.ErrNotFound) {
				w.WriteHeader(http.StatusNotFound)
				w.Write(apperror.ErrNotFound.Marshal())
				return
			}

			err = err.(*apperror.AppError)
			w.WriteHeader(http.StatusBadRequest)
			w.Write(appErr.Marshal())
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write(apperror.SystemError(err).Marshal())
	}
}
