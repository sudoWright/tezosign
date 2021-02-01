package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"msig/api/response"
	"msig/common/apperrors"
	"msig/common/log"
	"msig/models"
	"msig/repos"
	"msig/services"
	"msig/types"
	"net/http"
)

func (api *API) ContractStorageInit(w http.ResponseWriter, r *http.Request) {
	//TODO move to middleware
	net, err := ToNetwork(mux.Vars(r)["network"])
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

	service := services.New(repos.New(db), client, nil, net)

	resp, err := service.BuildContractInitStorage(req)
	if err != nil {
		//Unwrap apperror
		err, IsAppErr := apperrors.Unwrap(err)
		if !IsAppErr {
			log.Error("BuildContractInitStorage error: ", zap.Error(err))
		}

		response.JsonError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func (api *API) ContractStorageUpdate(w http.ResponseWriter, r *http.Request) {
	//todo add user

	//TODO move to middleware
	net, err := ToNetwork(mux.Vars(r)["network"])
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadParam, "network"))
		return
	}

	contractID := types.Address(mux.Vars(r)["contract_id"])
	if contractID == "" || contractID.Validate() != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadParam, "contract_id"))
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

	service := services.New(repos.New(db), client, nil, net)

	resp, err := service.BuildContractStorageUpdateOperation(contractID, req)
	if err != nil {
		//Unwrap apperror
		err, IsAppErr := apperrors.Unwrap(err)
		if !IsAppErr {
			log.Error("ContractStorageUpdate error: ", zap.Error(err))
		}

		response.JsonError(w, err)
		return
	}

	response.Json(w, resp)
}

func (api *API) ContractInfo(w http.ResponseWriter, r *http.Request) {
	//TODO move to middleware
	net, err := ToNetwork(mux.Vars(r)["network"])
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

	contractID := types.Address(mux.Vars(r)["contract_id"])
	if err = contractID.Validate(); err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadParam, "contract_id"))
		return
	}

	service := services.New(repos.New(db), client, nil, net)

	resp, err := service.ContractInfo(contractID)
	if err != nil {
		//Unwrap apperror
		err, IsAppErr := apperrors.Unwrap(err)
		if !IsAppErr {
			log.Error("BuildContractOperationToSign error: ", zap.Error(err))
		}

		response.JsonError(w, err)
		return
	}

	response.Json(w, resp)
}

func (api *API) ContractOperation(w http.ResponseWriter, r *http.Request) {
	//TODO move to middleware
	net, err := ToNetwork(mux.Vars(r)["network"])
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

	service := services.New(repos.New(db), client, nil, net)

	resp, err := service.ContractOperation(req)
	if err != nil {
		//Unwrap apperror
		err, IsAppErr := apperrors.Unwrap(err)
		if !IsAppErr {
			log.Error("BuildContractOperationToSign error: ", zap.Error(err))
		}

		response.JsonError(w, err)
		return
	}

	response.Json(w, resp)
}

func (api *API) OperationSignPayload(w http.ResponseWriter, r *http.Request) {
	user, isUser := r.Context().Value(ContextUserKey).(types.Address)
	if !isUser || (user.Validate() != nil) {
		response.JsonError(w, apperrors.New(apperrors.ErrService))
		return
	}

	net, err := ToNetwork(mux.Vars(r)["network"])
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadParam, "network"))
		return
	}

	operationID := mux.Vars(r)["operation_id"]
	if operationID == "" {
		response.JsonError(w, apperrors.New(apperrors.ErrBadParam, "operation_id"))
		return
	}

	payloadType := models.PayloadType(r.URL.Query().Get("type"))
	err = payloadType.Validate()
	if err != nil {
		response.JsonError(w, apperrors.NewWithDesc(apperrors.ErrBadParam, "type"))
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

	service := services.New(repos.New(db), client, nil, net)

	resp, err := service.BuildContractOperationToSign(user, operationID, payloadType)
	if err != nil {
		//Unwrap apperror
		err, IsAppErr := apperrors.Unwrap(err)
		if !IsAppErr {
			log.Error("BuildContractOperationReject error: ", zap.Error(err))
		}

		response.JsonError(w, err)
		return
	}

	response.Json(w, resp)
}

func (api *API) ContractOperationSignature(w http.ResponseWriter, r *http.Request) {
	user, isUser := r.Context().Value(ContextUserKey).(types.Address)
	if !isUser || (user.Validate() != nil) {
		response.JsonError(w, apperrors.New(apperrors.ErrService))
		return
	}

	//TODO move to middleware
	net, err := ToNetwork(mux.Vars(r)["network"])
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadParam, "network"))
		return
	}

	operationID := mux.Vars(r)["operation_id"]
	if operationID == "" {
		response.JsonError(w, apperrors.New(apperrors.ErrBadParam, "operation_id"))
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

	var req models.OperationSignature
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

	service := services.New(repos.New(db), client, nil, net)

	resp, err := service.SaveContractOperationSignature(user, operationID, req)
	if err != nil {
		//Unwrap apperror
		err, IsAppErr := apperrors.Unwrap(err)
		if !IsAppErr {
			log.Error("SaveContractOperationSignature error: ", zap.Error(err))
		}

		response.JsonError(w, err)
		return
	}

	response.Json(w, resp)
}

func (api *API) ContractOperationBuild(w http.ResponseWriter, r *http.Request) {
	user, isUser := r.Context().Value(ContextUserKey).(types.Address)
	if !isUser || (user.Validate() != nil) {
		response.JsonError(w, apperrors.New(apperrors.ErrService))
		return
	}

	//TODO move to middleware
	net, err := ToNetwork(mux.Vars(r)["network"])
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

	operationID, ok := mux.Vars(r)["operation_id"]
	if !ok || len(operationID) == 0 {
		response.JsonError(w, apperrors.NewWithDesc(apperrors.ErrBadParam, "tx_id"))
		return
	}

	payloadType := models.PayloadType(r.URL.Query().Get("type"))
	if err = payloadType.Validate(); err != nil {
		response.JsonError(w, apperrors.NewWithDesc(apperrors.ErrBadParam, "type"))
		return
	}

	service := services.New(repos.New(db), client, nil, net)

	resp, err := service.BuildContractOperation(user, operationID, payloadType)
	if err != nil {
		//Unwrap apperror
		err, IsAppErr := apperrors.Unwrap(err)
		if !IsAppErr {
			log.Error("ContractOperationBuild error: ", zap.Error(err))
		}

		response.JsonError(w, err)
		return
	}

	response.Json(w, resp)
}
