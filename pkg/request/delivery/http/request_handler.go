package http

import (
	"encoding/json"
	"net/http"

	"github.com/wascript3r/cryptopay/pkg/errcode"
	httpjson "github.com/wascript3r/httputil/json"

	"github.com/julienschmidt/httprouter"
	"github.com/wascript3r/anomaly/pkg/request"
)

type HTTPHandler struct {
	requestUcase request.Usecase
}

func NewHTTPHandler(r *httprouter.Router, ru request.Usecase) {
	handler := &HTTPHandler{
		requestUcase: ru,
	}

	r.POST("/api/request/process", handler.ProcessRequest)
	r.POST("/api/request/stats", handler.GetStats)
	r.POST("/api/request/all", handler.GetAll)
	r.POST("/api/request/imsi/stats", handler.GetIMSIStats)
	r.POST("/api/request/msc/stats", handler.GetMSCStats)
}

func serveError(w http.ResponseWriter, err error) {
	if err == request.InvalidInputError {
		httpjson.BadRequestCustom(w, request.InvalidInputError, nil)
		return
	}

	code := errcode.UnwrapErr(err, request.UnknownError)
	if code == request.UnknownError {
		httpjson.InternalErrorCustom(w, code, nil)
		return
	}

	httpjson.ServeErr(w, code, nil)
}

func (h *HTTPHandler) ProcessRequest(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req := &request.ProcessReq{}

	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		httpjson.BadRequest(w, nil)
		return
	}

	res, err := h.requestUcase.Process(r.Context(), req)
	if err != nil {
		serveError(w, err)
		return
	}

	httpjson.ServeJSON(w, res)
}

func (h *HTTPHandler) GetStats(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req := &request.FilterReq{}

	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		httpjson.BadRequest(w, nil)
		return
	}

	res, err := h.requestUcase.GetStats(r.Context(), req)
	if err != nil {
		serveError(w, err)
		return
	}

	httpjson.ServeJSON(w, res)
}

func (h *HTTPHandler) GetIMSIStats(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req := &request.AdvancedFilterReq{}

	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		httpjson.BadRequest(w, nil)
		return
	}

	res, err := h.requestUcase.GetIMSIStats(r.Context(), req)
	if err != nil {
		serveError(w, err)
		return
	}

	httpjson.ServeJSON(w, res)
}

func (h *HTTPHandler) GetMSCStats(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req := &request.AdvancedFilterReq{}

	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		httpjson.BadRequest(w, nil)
		return
	}

	res, err := h.requestUcase.GetMSCStats(r.Context(), req)
	if err != nil {
		serveError(w, err)
		return
	}

	httpjson.ServeJSON(w, res)
}

func (h *HTTPHandler) GetAll(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	req := &request.FilterReq{}

	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		httpjson.BadRequest(w, nil)
		return
	}

	res, err := h.requestUcase.GetAll(r.Context(), req)
	if err != nil {
		serveError(w, err)
		return
	}

	httpjson.ServeJSON(w, res)
}
