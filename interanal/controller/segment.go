package controller

import (
	"net/http"
	"strconv"
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
	if segment.Id == 0 {
		errorResponseJson(w, err.Error())
		return
	}

	if err != nil {
		logrus.Print(err)
		return
	}

	responseJson(w, segment)
}

type segmentCreateRequest struct {
	Slug string `json:"slug" validate:"required"`
}

func (h *segmentController) CreateSegment(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := r.ParseForm()
	if err != nil {
		logrus.Print(err)
		return
	}

	slug := r.PostFormValue("slug")
	logrus.Print(slug)
	var request = segmentCreateRequest{
		Slug: slug,
	}

	validate := validator.New()
	err = validate.Struct(request)
	if err != nil {
		errorResponseJson(w, err.Error())
		return
	}

	id, err := h.segmentService.CreateSegment(r.Context(), request.Slug)

	if id == 0 {
		errorResponseJson(w, err.Error())
		return
	}

	if err != nil {
		logrus.Print(err)
		return
	}

	type response struct { //todo В общий
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
		logrus.Print(err)
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
