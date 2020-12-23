package api

import (
	"msig/api/response"
	"msig/conf"
	"net/http"
)

func (api *API) Index(w http.ResponseWriter, r *http.Request) {
	response.Json(w, map[string]string{
		"service": conf.Service,
	})
}

func (api *API) Health(w http.ResponseWriter, r *http.Request) {
	response.Json(w, map[string]bool{
		"status": true,
	})
}
