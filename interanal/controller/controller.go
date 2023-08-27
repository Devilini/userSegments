package controller

import (
	"github.com/julienschmidt/httprouter"
)

type Handler struct {
}

//func (h *Handler) Register(router *httprouter.Router) {
//	router.HandlerFunc(http.MethodGet, "/api/user", h.GetUser)
//}

//func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
//	logrus.Print("TEST")
//	w.WriteHeader(204)
//	w.Write([]byte("TEST"))
//}

func (h *Handler) InitRoutes(router *httprouter.Router) {
	//router.HandlerFunc(http.MethodGet, "/api/user", GetUser())
}
