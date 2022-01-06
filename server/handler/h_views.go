package handler

import "net/http"

type viewsHandler struct{}

func NewViewsHandler() *viewsHandler {
	return &viewsHandler{}
}

func (f viewsHandler) Assets() http.Handler {
	assets := http.StripPrefix("/assets", http.FileServer(http.Dir("public/assets")))
	return assets
}

func (f viewsHandler) Views() http.Handler {
	views := http.FileServer(http.Dir("public/views"))
	return views
}
