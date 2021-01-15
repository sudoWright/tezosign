package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"msig/api/response"
	"msig/common/apperrors"
	"msig/common/log"
	models "msig/models"
	"msig/repos"
	"msig/services"
	"net/http"
)

func (api *API) ContractStorageInit(w http.ResponseWriter, r *http.Request) {
	//TODO move to middleware
	network, ok := mux.Vars(r)["network"]
	if !ok || network == "" {
		response.JsonError(w, apperrors.New(apperrors.ErrBadParam, "network"))
		return
	}

	net, err := ToNetwork(network)
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadParam, "network"))
		return
	}

	db, err := api.provider.GetDb(net)
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadParam))
		return
	}

	client, err := api.provider.GetRPCClient(net)
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadParam))
		return
	}

	var req models.ContractStorageRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadRequest))
		return
	}

	err = req.Validate()
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadRequest, err.Error()))
		return
	}

	service := services.New(repos.New(db), client, net)

	resp, err := service.BuildContractInitStorage(req)
	if err != nil {
		log.Error("BuildContractInitStorage error: ", zap.Error(err))
		response.JsonError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func (api *API) ContractStorageUpdate(w http.ResponseWriter, r *http.Request) {
	//TODO move to middleware
	network, ok := mux.Vars(r)["network"]
	if !ok || network == "" {
		response.JsonError(w, apperrors.New(apperrors.ErrBadParam, "network"))
		return
	}

	net, err := ToNetwork(network)
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadParam, "network"))
		return
	}

	db, err := api.provider.GetDb(net)
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadParam))
		return
	}

	client, err := api.provider.GetRPCClient(net)
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadParam))
		return
	}

	var req models.ContractStorageRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadRequest))
		return
	}

	err = req.Validate()
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadRequest, err.Error()))
		return
	}

	service := services.New(repos.New(db), client, net)

	resp, err := service.BuildContractStorageUpdateOperation(req)
	if err != nil {
		log.Error("BuildContractInitStorage error: ", zap.Error(err))
		response.JsonError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(resp))
}

func (api *API) ContractOperation(w http.ResponseWriter, r *http.Request) {
	//TODO move to middleware
	network, ok := mux.Vars(r)["network"]
	if !ok || network == "" {
		response.JsonError(w, apperrors.New(apperrors.ErrBadParam, "network"))
		return
	}

	net, err := ToNetwork(network)
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadParam, "network"))
		return
	}

	db, err := api.provider.GetDb(net)
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadParam))
		return
	}

	client, err := api.provider.GetRPCClient(net)
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadParam))
		return
	}

	var req models.ContractOperationRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadRequest))
		return
	}

	err = req.Validate()
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadRequest, err.Error()))
		return
	}

	service := services.New(repos.New(db), client, net)

	resp, err := service.BuildContractOperation(req)
	if err != nil {
		log.Error("BuildContractInitStorage error: ", zap.Error(err))
		response.JsonError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(resp))
}

func (api *API) ContractOperationSignature(w http.ResponseWriter, r *http.Request) {
	//Todo implement
}
