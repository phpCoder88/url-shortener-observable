package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/opentracing/opentracing-go"
)

// swagger:operation GET /report shortener URLReport
//
// Returns report information
//
// ---
// produces:
// - application/json
//
// parameters:
// - name: limit
//   in: query
//   description: Max number of records to return
//   required: false
//   type: integer
//   format: int64
//   default: 100
//
// - name: offset
//   in: query
//   description: Offset needed to return a specific subset of records
//   required: false
//   type: integer
//   format: int64
//   default: 0
//
// responses:
//   200:
//     description: returns report information
//     schema:
//       type: array
//       items:
//         "$ref": "#/definitions/ShortURLReportDto"
//   400:
//     description: Invalid input
//   500:
//     description: Internal error
func (h *Handler) ReportEndpoint(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("Content-Type", "application/json")

	span, ctx := opentracing.StartSpanFromContextWithTracer(req.Context(), h.tracer, "Handler.ReportEndpoint")
	defer span.Finish()

	var limit int64 = 100
	var offset int64 = 0
	var err error

	query := req.URL.Query()

	limit, err = h.container.ShortenerService.ParseLimitOffsetQueryParams(ctx, query, "limit", limit)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		h.logger.Error(err)
		return
	}

	offset, err = h.container.ShortenerService.ParseLimitOffsetQueryParams(ctx, query, "offset", offset)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		h.logger.Error(err)
		return
	}

	urls, err := h.container.ShortenerService.FindAll(ctx, limit, offset)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		h.logger.Error(err)
		return
	}

	if len(urls) > 0 {
		err = json.NewEncoder(res).Encode(urls)
	} else {
		_, err = res.Write([]byte(`[]`))
	}

	if err != nil {
		h.logger.Error(err)
		return
	}
}
