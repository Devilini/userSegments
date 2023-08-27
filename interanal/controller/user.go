package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"
	"userSegments/interanal/domain/user/service"
)

type userController struct {
	userService service.User
}

func NewUserController(productService service.User) *userController {
	r := &userController{
		userService: productService,
	}

	//g.POST("/create", r.create)
	//g.GET("/", r.getById)

	return r
}

func (h *userController) GetUser(w http.ResponseWriter, r *http.Request) {
	has := r.URL.Query().Has("id")
	logrus.Print(has)
	if !has {
		errorResponseJson(w, "Err, need user id")
		return
	}
	id := r.URL.Query().Get("id")
	// string to int
	idInt, err := strconv.Atoi(id)
	if err != nil {
		panic(err)
	}
	user, err := h.userService.GetUserById(r.Context(), idInt)
	if err != nil {
		logrus.Print("err")
		logrus.Print(err)
		//return err
	}
	logrus.Print(user)

	responseJson(w, user)
}

func responseJson(w http.ResponseWriter, data any) {
	jsonResp, err := json.Marshal(data)
	if err != nil {
		// handle error
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}

func errorResponseJson(w http.ResponseWriter, text string) {
	resp := make(map[string]string)
	resp["error"] = text
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	w.Write(jsonResp)

	//err := errors.New(text)
	//w.Header().Set("Content-Type", "application/json")
	//w.WriteHeader(http.StatusBadRequest)
	//
	//w.Write([]byte(err.Error()))
}
