package controller

import (
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"userSegments/interanal/domain/user/service"
)

type userController struct {
	userService service.User
}

func NewUserController(productService service.User) *userController {
	return &userController{
		userService: productService,
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
		logrus.Print(err)
		return
	}

	responseJson(w, user)
}

type userCreateRequest struct {
	Name string `json:"name" validate:"required"`
}

func (h *userController) CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := r.ParseForm()
	if err != nil {
		logrus.Print(err)
		return
	}

	userName := r.PostFormValue("name")
	logrus.Print(userName)
	//logrus.Print(name)
	//return
	//if userName == "" {
	//	logrus.Print("Empty user name")
	//	return
	//}
	var request = userCreateRequest{
		Name: userName,
	}

	//resp := make(map[string]string)
	//resp["name"] = "text"
	//_ = json.Unmarshal([]byte(resp), &u)
	//logrus.Print(request)
	validate := validator.New()
	err = validate.Struct(request)
	if err != nil {
		errorResponseJson(w, err.Error())
		return
	}

	id, err := h.userService.CreateUser(r.Context(), request.Name)
	if err != nil {
		logrus.Print(err)
		return
	}

	type response struct {
		Id int `json:"id"`
	}

	responseJson(w, response{
		Id: id,
	})
}
