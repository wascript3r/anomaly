package http

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/wascript3r/anomaly/pkg/graph"
	"github.com/wascript3r/cryptopay/pkg/errcode"
	httpjson "github.com/wascript3r/httputil/json"
)

type HTTPHandler struct {
	graphUcase graph.Usecase
}

func NewHTTPHandler(r *httprouter.Router, gu graph.Usecase) {
	handler := &HTTPHandler{
		graphUcase: gu,
	}

	r.GET("/api/graph/all", handler.AllGraphs)
}

func serveError(w http.ResponseWriter, err error) {
	if err == graph.InvalidInputError {
		httpjson.BadRequestCustom(w, graph.InvalidInputError, nil)
		return
	}

	code := errcode.UnwrapErr(err, graph.UnknownError)
	if code == graph.UnknownError {
		httpjson.InternalErrorCustom(w, code, nil)
		return
	}

	httpjson.ServeErr(w, code, nil)
}

func (h *HTTPHandler) AllGraphs(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	res, err := h.graphUcase.GetAll(r.Context())
	if err != nil {
		serveError(w, err)
		return
	}

	httpjson.ServeJSON(w, res)
}
