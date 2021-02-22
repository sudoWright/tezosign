package api

import (
	"net/http"
	"tezosign/api/response"
	"tezosign/repos"
	"tezosign/services"
)

func (api *API) TezosExchangeRates(w http.ResponseWriter, r *http.Request) {
	net, networkContext, err := GetNetworkContext(r)
	if err != nil {
		response.JsonError(w, err)
		return
	}

	service := services.New(repos.New(networkContext.Db), repos.New(networkContext.IndexerDB), networkContext.Client, networkContext.Auth, net)

	rates, err := service.TezosExchangeRates()
	if err != nil {
		response.JsonError(w, err)
		return
	}

	response.Json(w, rates)
}
