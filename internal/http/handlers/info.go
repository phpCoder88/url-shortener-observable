package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/phpCoder88/url-shortener/internal/responses"
	"github.com/phpCoder88/url-shortener/internal/version"
)

// swagger:route GET /service-info service BuildInfo
//
// Returns build information
//
// Produces:
// - application/json
//
// Responses:
//   200: BuiltInfo
func (h *Handler) BuiltInfoEndpoint(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("Content-Type", "application/json")

	buildInfo := responses.BuiltInfo{
		Version:     version.Version,
		BuildDate:   version.BuildDate,
		BuildCommit: version.BuildCommit,
	}

	err := json.NewEncoder(res).Encode(buildInfo)
	if err != nil {
		h.logger.Error(err)
	}
}
