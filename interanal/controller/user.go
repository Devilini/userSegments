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

type userController struct {
	userService service.User
}

func NewUserController(userService service.User) *userController {
	return &userController{
		userService: userService,
	}
}

func (h *userController) GetUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		errorResponseJson(w, "Invalid format of user id")
		return
	}

	user, err := h.userService.GetUserById(r.Context(), id)
	if err != nil {
		logrus.Info(err.Error())
		errorResponseJson(w, "User does not exists")
		return
	}

	responseJson(w, user)
}

func (h *userController) CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	var req request.UserCreateRequest
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

	id, err := h.userService.CreateUser(r.Context(), req.Name)
	if err != nil {
		logrus.Error(err.Error())
		errorResponseJson(w, "Cannot create user")
		return
	}

	type response struct {
		Id int `json:"id"`
	}

	responseJson(w, response{
		Id: id,
	})
}
