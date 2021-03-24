package api

import (
	"encoding/json"
	"tezosign/common/log"

	"go.uber.org/zap"

	"net/http"
	"tezosign/api/response"
	"tezosign/common/apperrors"
	"tezosign/models"
	"tezosign/repos"
	"tezosign/services"
	"tezosign/types"

	"github.com/gorilla/mux"
)

func (api *API) AssetsList(w http.ResponseWriter, r *http.Request) {
	user, net, networkContext, err := GetUserNetworkContext(r)
	if err != nil {
		response.JsonError(w, err)
		return
	}

	contractAddress := types.Address(mux.Vars(r)[ContractIDParam])
	err = contractAddress.Validate()
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadParam, ContractIDParam))
		return
	}

	service := services.New(repos.New(networkContext.Db), repos.New(networkContext.IndexerDB), networkContext.Client, networkContext.Auth, net)

	assets, err := service.AssetsList(user, contractAddress)
	if err != nil {
		//Unwrap apperror
		err, IsAppErr := apperrors.Unwrap(err)
		if !IsAppErr {
			log.Error("AssetsList error: ", zap.Error(err))
		}

		response.JsonError(w, err)
		return
	}

	response.Json(w, assets)
}

func (api *API) AssetsExchangeRates(w http.ResponseWriter, r *http.Request) {
	user, net, networkContext, err := GetUserNetworkContext(r)
	if err != nil {
		response.JsonError(w, err)
		return
	}

	contractAddress := types.Address(mux.Vars(r)[ContractIDParam])
	err = contractAddress.Validate()
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadParam, ContractIDParam))
		return
	}

	service := services.New(repos.New(networkContext.Db), repos.New(networkContext.IndexerDB), networkContext.Client, networkContext.Auth, net)

	assets, err := service.AssetsExchangeRates(user, contractAddress)
	if err != nil {
		response.JsonError(w, err)
		return
	}

	response.Json(w, assets)
}

func (api *API) ContractAsset(w http.ResponseWriter, r *http.Request) {
	user, net, networkContext, err := GetUserNetworkContext(r)
	if err != nil {
		response.JsonError(w, err)
		return
	}

	contractAddress := types.Address(mux.Vars(r)[ContractIDParam])
	err = contractAddress.Validate()
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadParam, ContractIDParam))
		return
	}

	var data models.Asset
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadRequest))
		return
	}

	if err = data.Validate(); err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadRequest, err.Error()))
		return
	}

	service := services.New(repos.New(networkContext.Db), repos.New(networkContext.IndexerDB), networkContext.Client, networkContext.Auth, net)

	reps, err := service.ContractAsset(user, contractAddress, data)
	if err != nil {
		//Unwrap apperror
		err, IsAppErr := apperrors.Unwrap(err)
		if !IsAppErr {
			log.Error("ContractAsset error: ", zap.Error(err))
		}

		response.JsonError(w, err)
		return
	}

	response.Json(w, reps)
}

func (api *API) ContractAssetEdit(w http.ResponseWriter, r *http.Request) {
	user, net, networkContext, err := GetUserNetworkContext(r)
	if err != nil {
		response.JsonError(w, err)
		return
	}

	contractAddress := types.Address(mux.Vars(r)[ContractIDParam])
	err = contractAddress.Validate()
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadParam, ContractIDParam))
		return
	}

	var data models.Asset
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadRequest))
		return
	}

	if err = data.Validate(); err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadRequest, err.Error()))
		return
	}

	service := services.New(repos.New(networkContext.Db), repos.New(networkContext.IndexerDB), networkContext.Client, networkContext.Auth, net)

	reps, err := service.ContractAssetEdit(user, contractAddress, data)
	if err != nil {
		//Unwrap apperror
		err, IsAppErr := apperrors.Unwrap(err)
		if !IsAppErr {
			log.Error("ContractAssetEdit error: ", zap.Error(err))
		}

		response.JsonError(w, err)
		return
	}

	response.Json(w, reps)
}

func (api *API) RemoveContractAsset(w http.ResponseWriter, r *http.Request) {
	user, net, networkContext, err := GetUserNetworkContext(r)
	if err != nil {
		response.JsonError(w, err)
		return
	}

	contractAddress := types.Address(mux.Vars(r)[ContractIDParam])
	err = contractAddress.Validate()
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadParam, ContractIDParam))
		return
	}

	var data models.Asset
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadRequest))
		return
	}

	if err = data.Address.Validate(); err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadParam, "address"))
		return
	}

	service := services.New(repos.New(networkContext.Db), repos.New(networkContext.IndexerDB), networkContext.Client, networkContext.Auth, net)

	err = service.RemoveContractAsset(user, contractAddress, data)
	if err != nil {
		//Unwrap apperror
		err, IsAppErr := apperrors.Unwrap(err)
		if !IsAppErr {
			log.Error("RemoveContractAsset error: ", zap.Error(err))
		}

		response.JsonError(w, err)
		return
	}

	response.Json(w, map[string]interface{}{"message": "success"})
}
