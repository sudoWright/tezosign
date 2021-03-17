package api

import (
	"encoding/json"
	"net/http"
	"tezosign/api/response"
	"tezosign/common/apperrors"
	"tezosign/common/log"
	"tezosign/models"
	"tezosign/repos"
	"tezosign/services"

	"go.uber.org/zap"
)

func (api *API) VestingContractStorageInit(w http.ResponseWriter, r *http.Request) {
	//Use GetUserNetworkContext to check user middleware
	_, net, networkContext, err := GetUserNetworkContext(r)
	if err != nil {
		response.JsonError(w, err)
		return
	}

	var req models.VestingContractStorageRequest
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

	service := services.New(repos.New(networkContext.Db), repos.New(networkContext.IndexerDB), networkContext.Client, networkContext.Auth, net)

	resp, err := service.BuildVestingContractInitStorage(req)
	if err != nil {
		//Unwrap apperror
		err, IsAppErr := apperrors.Unwrap(err)
		if !IsAppErr {
			log.Error("BuildVestingContractInitStorage error: ", zap.Error(err))
		}

		response.JsonError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}
