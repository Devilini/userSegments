package controller

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"userSegments/interanal/controller/request"
	"userSegments/interanal/model"
	"userSegments/interanal/service"
	"userSegments/pkg/helper"
)

type userSegmentsController struct {
	userSegmentsService service.UserSegments
}

func NewUserSegmentsController(userSegmentsService service.UserSegments) *userSegmentsController {
	return &userSegmentsController{
		userSegmentsService: userSegmentsService,
	}
}

func (h *userSegmentsController) GetUserSegments(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		errorResponseJson(w, "Invalid format of user id")
		return
	}
	segments, err := h.userSegmentsService.GetUserSegments(r.Context(), id)
	if err != nil {
		logrus.Info(err)
		errorResponseJson(w, err.Error())
		return
	}

	type response struct {
		Segments []model.Segment `json:"segments"`
		UserId   int             `json:"user_id"`
	}
	responseJson(w, response{
		Segments: segments,
		UserId:   id,
	})
}

func (h *userSegmentsController) ChangeUserSegments(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		errorResponseJson(w, "Invalid format of user id")
		return
	}

	decoder := json.NewDecoder(r.Body)
	var req request.UserAddSegmentRequest
	if err := decoder.Decode(&req); err != nil {
		logrus.Info(err)
		errorResponseJson(w, "Error parse params")
		return
	}
	req.UserId = id

	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		errorResponseJson(w, err.Error())
		return
	}

	if len(req.AddSegments) == 0 && len(req.DeleteSegments) == 0 {
		errorResponseJson(w, "All segment structs are empty")
		return
	}

	if len(helper.IntersectionSlices(req.AddSegments, req.DeleteSegments)) > 0 {
		errorResponseJson(w, "Segment structs has duplicates")
		return
	}

	result, err := h.userSegmentsService.ChangeUserSegments(r.Context(), req)
	if err != nil {
		logrus.Info(err)
		errorResponseJson(w, "Error change segments")
		return
	}

	type response struct {
		Status string `json:"status"`
	}

	status := "Failed"
	if result > 0 {
		status = "Success"
	}

	responseJson(w, response{
		Status: status,
	})
}
