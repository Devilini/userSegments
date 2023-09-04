package controller

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
	"userSegments/interanal/apperror"
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
		errorResponseJson(w, apperror.NewAppError(err, "Invalid format of user id"))
		return
	}
	segments, err := h.userSegmentsService.GetUserSegments(r.Context(), id)
	if err != nil {
		errorResponseJson(w, err)
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
		errorResponseJson(w, apperror.NewAppError(err, "Invalid format of user id"))
		return
	}

	decoder := json.NewDecoder(r.Body)
	var req request.UserAddSegmentRequest
	if err := decoder.Decode(&req); err != nil {
		errorResponseJson(w, apperror.NewAppError(err, "Error parse params"))
		return
	}
	req.UserId = id

	validate := validator.New()
	err = validate.Struct(req)
	if err != nil {
		errorResponseJson(w, apperror.NewAppError(err, err.Error()))
		return
	}

	if len(req.AddSegments) == 0 && len(req.DeleteSegments) == 0 {
		errorResponseJson(w, apperror.NewAppError(err, "All segment structs are empty"))
		return
	}

	if len(helper.IntersectionSlices(req.AddSegments, req.DeleteSegments)) > 0 {
		errorResponseJson(w, apperror.NewAppError(err, "Segment structs has duplicates"))
		return
	}

	result, err := h.userSegmentsService.ChangeUserSegments(r.Context(), req)
	if err != nil {
		errorResponseJson(w, err)
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
