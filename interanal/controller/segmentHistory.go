package controller

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"log"
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

	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))

	http.ServeFile(w, r, path)

	file.Close()
}

func (h *segmentHistoryController) GenerateHistoryReport(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var req struct {
		Date string `json:"date" validate:"required"`
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		logrus.Print(err)
		return
	}
	log.Println(req)
	path, err := h.segmentHistoryService.GetSegmentsReport(r.Context(), req.Date)
	if err != nil {
		errorResponseJson(w, err.Error())
		return
	}

	type response struct {
		Link string `json:"link"`
	}
	responseJson(w, response{
		Link: fmt.Sprintf("/reports/segments_%s", path),
	})
}
