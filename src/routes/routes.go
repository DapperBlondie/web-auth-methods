package routes

import (
	"github.com/DapperBlondie/web-auth-methods/src/handlers"
	"github.com/go-chi/chi"
	"net/http"
)

func AppRoutes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(handlers.Conf.EnableCORS)
	mux.Use(handlers.Conf.LoadSession)
	mux.Get("/status", handlers.Conf.CheckStatusHandler)

	return mux
}
