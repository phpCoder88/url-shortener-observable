package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"

	"github.com/phpCoder88/url-shortener-observable/internal/helpers"
)

func (h *Handler) RedirectFullURL(res http.ResponseWriter, req *http.Request) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(req.Context(), h.tracer, "Handler.RedirectFullURL")
	defer span.Finish()

	token := mux.Vars(req)["token"]
	h.logger.Infof("Requested token %s", mux.Vars(req)["token"])

	userIP, err := helpers.GetIP(ctx, req)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		h.logger.Error(err)
		return
	}

	fullURL, err := h.container.ShortenerService.VisitFullURL(ctx, token, userIP)
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		h.logger.Error(err)
		return
	}

	http.Redirect(res, req, fullURL, http.StatusSeeOther)
}
