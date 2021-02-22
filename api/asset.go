package api

import (
	"log"
	"net/http"
	"tezosign/api/response"
	"tezosign/repos"
	"tezosign/services"
)

func (api *API) AssetsList(w http.ResponseWriter, r *http.Request) {
	net, networkContext, err := GetNetworkContext(r)
	if err != nil {
		response.JsonError(w, err)
		return
	}

	service := services.New(repos.New(networkContext.Db), repos.New(networkContext.IndexerDB), networkContext.Client, networkContext.Auth, net)

	assets, err := service.AssetsList()
	if err != nil {
		response.JsonError(w, err)
		return
	}

	response.Json(w, assets)
}

func (api *API) AssetsExchangeRates(w http.ResponseWriter, r *http.Request) {
	net, networkContext, err := GetNetworkContext(r)
	if err != nil {
		response.JsonError(w, err)
		return
	}

	service := services.New(repos.New(networkContext.Db), repos.New(networkContext.IndexerDB), networkContext.Client, networkContext.Auth, net)

	assets, err := service.AssetsExchangeRates()
	if err != nil {
		log.Print(err)
		response.JsonError(w, err)
		return
	}

	response.Json(w, assets)
}
