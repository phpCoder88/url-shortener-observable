package handlers

import (
	"net/http"

	"github.com/phpCoder88/url-shortener/internal/helpers"

	"github.com/gorilla/mux"
)

func (h *Handler) RedirectFullURL(res http.ResponseWriter, req *http.Request) {
	token := mux.Vars(req)["token"]
	h.logger.Infof("Requested token %s", mux.Vars(req)["token"])

	userIP, err := helpers.GetIP(req)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		h.logger.Error(err)
		return
	}

	fullURL, err := h.container.ShortenerService.VisitFullURL(token, userIP)
	if err != nil {
		res.WriteHeader(http.StatusNotFound)
		h.logger.Error(err)
		return
	}

	http.Redirect(res, req, fullURL, http.StatusSeeOther)
}
