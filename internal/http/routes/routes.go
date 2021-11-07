package routes

import (
	"net/http"

	gHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"

	"github.com/phpCoder88/url-shortener-observable/internal/http/handlers"
	"github.com/phpCoder88/url-shortener-observable/internal/http/middlewares"
	"github.com/phpCoder88/url-shortener-observable/internal/ioc"
)

func Routes(logger *zap.SugaredLogger, container *ioc.Container, tracer opentracing.Tracer) http.Handler {
	standardMiddleware := alice.New(middlewares.RecoverPanic(logger))

	router := mux.NewRouter()
	api := router.PathPrefix("/api").Subrouter()

	handler := handlers.NewHandler(logger, container, tracer)

	api.HandleFunc("/shorten", handler.ShortenEndpoint).Methods("POST")
	api.HandleFunc("/report", handler.ReportEndpoint).Methods("GET")
	api.HandleFunc("/service-info", handler.BuiltInfoEndpoint).Methods("GET")
	router.PathPrefix("/swaggerui/").Handler(http.StripPrefix("/swaggerui/", http.FileServer(http.Dir("./web/static/swaggerui"))))
	router.HandleFunc("/", handler.RedirectFullURL).Methods("GET").Queries("t", "{token}")

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
