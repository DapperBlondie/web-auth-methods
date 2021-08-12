package routes

import (
	"github.com/DapperBlondie/web-auth-methods/src/handlers"
	"github.com/go-chi/chi"
	"net/http"
)

// AppRoutes use for creating and managing routes of our application
func AppRoutes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(handlers.Conf.EnableCORS)
	mux.Use(handlers.Conf.LoadSession)
	mux.Get("/status", handlers.Conf.CheckStatusHandler)
	mux.Get("/save-hmac-token", handlers.Conf.SaveHmacToken)
	mux.Get("/get-hmac-token", handlers.Conf.GetAndCheckHmacToken)

	return mux
}
