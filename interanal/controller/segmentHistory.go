package controller

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"userSegments/interanal/service"
)

type segmentHistoryController struct {
	segmentHistoryService service.SegmentHistory
}

func NewSegmentHistoryControllerController(segmentHistoryService service.SegmentHistory) *segmentHistoryController {
	return &segmentHistoryController{
		segmentHistoryService: segmentHistoryService,
	}
}

func (h *segmentHistoryController) DownloadReport(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fileName := ps.ByName("filename")
	path := fmt.Sprintf("public/reports/%s", fileName)

	file, err := os.Open(path)
	if err != nil && os.IsNotExist(err) {
		errorResponseJson(w, "report not found")
		return
	}
	defer file.Close()

	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))

	http.ServeFile(w, r, path)
}

func (h *segmentHistoryController) GenerateHistoryReport(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var req struct {
		DateFrom string `json:"dateFrom" validate:"required"`
		DateTo   string `json:"dateTo" validate:"required"`
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		logrus.Error(err)
		errorResponseJson(w, "Error parse params")
		return
	}
	fileName, err := h.segmentHistoryService.GetSegmentsReport(r.Context(), req.DateFrom, req.DateTo)
	if err != nil {
		logrus.Error(err)
		errorResponseJson(w, "Report generating error")
		return
	}

	type response struct {
		Link string `json:"link"`
	}
	responseJson(w, response{
		Link: fmt.Sprintf("/reports/segments_%s", fileName),
	})
}
