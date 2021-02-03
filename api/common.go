package api

import (
	"net/http"
	"tezosign/api/response"
	"tezosign/conf"
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
