package metric

import (
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Handler struct {
}

func (h *Handler) Register(router *httprouter.Router) { //todo удалить
	router.HandlerFunc(http.MethodGet, "/api/heartbeat", h.Heartbeat)
}

func (h *Handler) Heartbeat(w http.ResponseWriter, r *http.Request) {
	logrus.Print("TEST")
	w.WriteHeader(204)
	w.Write([]byte("TEST"))
}
