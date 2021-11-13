package routes

import (
	"net/http"

	"github.com/phpCoder88/url-shortener-observable/internal/http/handlers"
	"github.com/phpCoder88/url-shortener-observable/internal/http/middlewares"
	"github.com/phpCoder88/url-shortener-observable/internal/ioc"

	gHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"go.uber.org/zap"
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
