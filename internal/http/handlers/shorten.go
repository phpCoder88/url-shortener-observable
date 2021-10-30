package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"gopkg.in/go-playground/validator.v9"

	"github.com/phpCoder88/url-shortener/internal/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

// swagger:operation POST /shorten shortener shortenURL
//
// Creates a new short URL for given URL
//
// ---
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// parameters:
//   - in: body
//     name: url
//     description: URL to shorten
//     schema:
//       type: object
//       required:
//         - fullURL
//       properties:
//         fullURL:
//           type: string
//
// Responses:
//   200:
//     description: Short URL already exists
//     schema:
//       type: object
//       properties:
//         shortURL:
//           type: string
//           description: Short URL
//   201:
//     description: Created short URL
//     schema:
//       type: object
//       properties:
//         shortURL:
//           type: string
//           description: Short URL
//   400:
//     description: Invalid input
//   500:
//     description: Internal error
func (h *Handler) ShortenEndpoint(res http.ResponseWriter, req *http.Request) {
	start := time.Now()
	defer func() {
		h.metrics.LatencyHistogram.
			With(prometheus.Labels{metrics.LabelHandler: "ShortenEndpoint"}).
			Observe(float64(time.Since(start).Milliseconds()))
	}()

	res.Header().Add("Content-Type", "application/json")

	var input struct {
		FullURL string `validate:"required,url"`
	}

	err := json.NewDecoder(req.Body).Decode(&input)
	if err != nil {
		h.logger.Error(err)
	}

	validate := validator.New()
	err = validate.Struct(input)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		h.logger.Error(err)
		return
	}

	h.logger.Info(input)
	urlModel, exist, err := h.container.ShortenerService.CreateShortURL(input.FullURL)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		h.logger.Error(err)
		return
	}

	if exist {
		res.WriteHeader(http.StatusOK)
	} else {
		res.WriteHeader(http.StatusCreated)
	}

	shortURL := fmt.Sprintf("http://%s?t=%s", req.Host, urlModel.Token)
	h.logger.Info(shortURL)
	_, err = res.Write([]byte(`{"shortURL": "` + shortURL + `"}`))
	if err != nil {
		h.logger.Error(err)
	}
}
