package controller

import (
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"userSegments/interanal/model"
	"userSegments/interanal/service"
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
		logrus.Print(err)
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

//type userCreateRequest struct {
//	Name string `json:"name" validate:"required"`
//}
//
//func (h *userController) CreateUser2(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
//	err := r.ParseForm()
//	if err != nil {
//		logrus.Print(err)
//		return
//	}
//
//	userName := r.PostFormValue("name")
//
//	var request = userCreateRequest{
//		Name: userName,
//	}
//
//	validate := validator.New()
//	err = validate.Struct(request)
//	if err != nil {
//		errorResponseJson(w, err.Error())
//		return
//	}
//
//	id, err := h.userService.CreateUser(r.Context(), request.Name)
//	if err != nil {
//		logrus.Print(err)
//		return
//	}
//
//	type response struct {
//		Id int `json:"id"`
//	}
//
//	responseJson(w, response{
//		Id: id,
//	})
//}
