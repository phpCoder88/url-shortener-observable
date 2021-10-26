package routes

import (
	"net/http"

	gHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"github.com/phpCoder88/url-shortener/internal/http/handlers"
	"github.com/phpCoder88/url-shortener/internal/http/middlewares"
	"github.com/phpCoder88/url-shortener/internal/ioc"
)

func Routes(logger *zap.SugaredLogger, container *ioc.Container) http.Handler {
	standardMiddleware := alice.New(middlewares.RecoverPanic(logger))

	router := mux.NewRouter()
	api := router.PathPrefix("/api").Subrouter()

	handler := handlers.NewHandler(logger, container)

	api.HandleFunc("/shorten", handler.ShortenEndpoint).Methods("POST")
	api.HandleFunc("/report", handler.ReportEndpoint).Methods("GET")
	api.HandleFunc("/service-info", handler.BuiltInfoEndpoint).Methods("GET")
	router.PathPrefix("/swaggerui/").Handler(http.StripPrefix("/swaggerui/", http.FileServer(http.Dir("./web/static/swaggerui"))))
	router.HandleFunc("/", handler.RedirectFullURL).Methods("GET").Queries("t", "{token}")
	router.Handle("/metrics", promhttp.Handler()).Methods("GET")

	methods := gHandlers.AllowedMethods([]string{
		"GET",
		"POST",
	})
	headers := gHandlers.AllowedHeaders([]string{
		"Content-Type",
		"Authorization",
		"X-Requested-With",
	})
	origins := gHandlers.AllowedOrigins([]string{"*"})

	return gHandlers.CORS(headers, methods, origins)(standardMiddleware.Then(router))
}
