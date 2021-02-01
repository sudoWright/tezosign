package api

import (
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"msig/api/response"
	"msig/common/apperrors"
	"msig/common/log"
	"msig/repos"
	"msig/services"
	"msig/types"
	"net/http"
)

func (api *API) ContractOperationsList(w http.ResponseWriter, r *http.Request) {
	user, isUser := r.Context().Value(ContextUserKey).(types.Address)
	if !isUser || (user.Validate() != nil) {
		response.JsonError(w, apperrors.New(apperrors.ErrService))
		return
	}

	vars := mux.Vars(r)

	net, err := ToNetwork(vars["network"])
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadParam, "network"))
		return
	}

	contractAddress := types.Address(vars["contract_id"])
	err = contractAddress.Validate()
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadParam, "contract_id"))
		return
	}

	//TODO add params
	//params := map[string]interface{}{}
	//err = api.queryDecoder.Decode(&params, r.URL.Query())
	//if err != nil {
	//	response.JsonError(w, apperrors.New(apperrors.ErrBadRequest))
	//	return
	//}

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
