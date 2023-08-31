package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"userSegments/interanal/controller/request"
	"userSegments/interanal/service"

	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

type segmentController struct {
	segmentService service.Segment
}

func NewSegmentController(segmentService service.Segment) *segmentController {
	return &segmentController{
		segmentService: segmentService,
	}
}

func (h *segmentController) GetSegment(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		errorResponseJson(w, "Invalid format of segment id")
		return
	}

	segment, err := h.segmentService.GetSegmentById(r.Context(), id)
	if err != nil {
		logrus.Info(err)
		errorResponseJson(w, "Segment does not exists")
		return
	}

	responseJson(w, segment)
}

func (h *segmentController) CreateSegment(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	var req request.SegmentCreateRequest
	if err := decoder.Decode(&req); err != nil {
		logrus.Error(err)
		errorResponseJson(w, "Error parse params")
		return
	}

	validate := validator.New()
	err := validate.Struct(req)
	if err != nil {
		errorResponseJson(w, err.Error())
		return
	}

	id, err := h.segmentService.CreateSegment(r.Context(), req.Slug)
	if id == 0 {
		errorResponseJson(w, err.Error())
		return
	}

	if err != nil {
		logrus.Print(err)
		errorResponseJson(w, "Cannot create segment")
		return
	}

	type response struct {
		Id int `json:"id"`
	}

	responseJson(w, response{
		Id: id,
	})
}

func (h *segmentController) DeleteSegment(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	slug := ps.ByName("slug")
	err := h.segmentService.DeleteSegmentBySlug(r.Context(), slug)
	if err != nil {
		logrus.Info(err)
		errorResponseJson(w, err.Error())
		return
	}

	type response struct {
		Status string `json:"status"`
	}

	responseJson(w, response{
		Status: "Success",
	})
}
