package api

import (
	"net/http"
	"tezosign/api/response"
	"tezosign/common/apperrors"
	"tezosign/common/log"
	"tezosign/repos"
	"tezosign/services"
	"tezosign/types"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func (api *API) ContractOperationsList(w http.ResponseWriter, r *http.Request) {
	user, net, networkContext, err := GetUserNetworkContext(r)
	if err != nil {
		response.JsonError(w, err)
		return
	}

	contractAddress := types.Address(mux.Vars(r)["contract_id"])
	err = contractAddress.Validate()
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadParam, "contract_id"))
		return
	}

	service := services.New(repos.New(networkContext.Db), repos.New(networkContext.IndexerDB), networkContext.Client, nil, net)

	list, err := service.GetOperationsList(user, contractAddress, nil)
	if err != nil {
		//Unwrap apperror
		err, IsAppErr := apperrors.Unwrap(err)
		if !IsAppErr {
			log.Error("ContractOperationsList error: ", zap.Error(err))
		}

		response.JsonError(w, err)
		return
	}

	response.Json(w, list)
}

func (api *API) ContractOrigination(w http.ResponseWriter, r *http.Request) {
	net, networkContext, err := GetNetworkContext(r)
	if err != nil {
		response.JsonError(w, err)
		return
	}

	txID := mux.Vars(r)["tx_id"]
	if txID == "" {
		response.JsonError(w, apperrors.New(apperrors.ErrBadParam, "tx_id"))
		return
	}

	service := services.New(repos.New(networkContext.Db), repos.New(networkContext.IndexerDB), networkContext.Client, networkContext.Auth, net)

	contractID, err := service.CheckContractOrigination(txID)
	if err != nil {
		response.JsonError(w, err)
		return
	}

	response.Json(w, map[string]interface{}{
		"contract": contractID,
	})
}
