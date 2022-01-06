package handler

import (
	"net/http"

	"github.com/fahmifj/ag-portal/service"
	"github.com/gorilla/mux"
)

type Router struct {
	service *service.Service
	*mux.Router
}

func NewRouter(service *service.Service) *Router {
	router := &Router{
		Router:  mux.NewRouter(),
		service: service,
	}

	vmFetchHandler := router.Methods(http.MethodGet).Subrouter()
	vmFetchHandler.HandleFunc("/vm/fetch", router.FetchVM) // should be get list
	vmCtlHandler := router.Methods(http.MethodPost).Subrouter()
	vmCtlHandler.HandleFunc("/vm/{vm:[a-zA-Z0-9\\-]+}", router.VMControlHandler).Queries("action", "{start|stop}")

	fsHandler := NewViewsHandler()
	router.Handle("/", fsHandler.Views())
	router.PathPrefix("/assets").Handler(fsHandler.Assets())

	return router
}
