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
	"tezosign/types"

	"github.com/gorilla/mux"

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

func (api *API) VestingContractOperation(w http.ResponseWriter, r *http.Request) {
	_, net, networkContext, err := GetUserNetworkContext(r)
	if err != nil {
		response.JsonError(w, err)
		return
	}

	var req models.VestingContractOperation
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

	resp, err := service.VestingContractOperation(req)
	if err != nil {
		//Unwrap apperror
		err, IsAppErr := apperrors.Unwrap(err)
		if !IsAppErr {
			log.Error("VestingContractOperation error: ", zap.Error(err))
		}

		response.JsonError(w, err)
		return
	}

	response.Json(w, resp)
}

func (api *API) VestingContractInfo(w http.ResponseWriter, r *http.Request) {
	//Use GetUserNetworkContext to check user middleware
	_, net, networkContext, err := GetUserNetworkContext(r)
	if err != nil {
		response.JsonError(w, err)
		return
	}

	contractID := types.Address(mux.Vars(r)["contract_id"])
	if err := contractID.Validate(); err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadParam, "contract_id"))
		return
	}

	service := services.New(repos.New(networkContext.Db), repos.New(networkContext.IndexerDB), networkContext.Client, networkContext.Auth, net)

	resp, err := service.VestingContractInfo(contractID)
	if err != nil {
		//Unwrap apperror
		err, IsAppErr := apperrors.Unwrap(err)
		if !IsAppErr {
			log.Error("VestingContractInfo error: ", zap.Error(err))
		}

		response.JsonError(w, err)
		return
	}

	response.Json(w, resp)
}

func (api *API) VestingsList(w http.ResponseWriter, r *http.Request) {
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

	reps, err := service.VestingsList(user, contractAddress)
	if err != nil {
		//Unwrap apperror
		err, IsAppErr := apperrors.Unwrap(err)
		if !IsAppErr {
			log.Error("VestingsList error: ", zap.Error(err))
		}

		response.JsonError(w, err)
		return
	}

	response.Json(w, reps)
}

func (api *API) ContractVesting(w http.ResponseWriter, r *http.Request) {
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

	var data models.Vesting
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

	reps, err := service.ContractVesting(user, contractAddress, data)
	if err != nil {
		//Unwrap apperror
		err, IsAppErr := apperrors.Unwrap(err)
		if !IsAppErr {
			log.Error("ContractVesting error: ", zap.Error(err))
		}

		response.JsonError(w, err)
		return
	}

	response.Json(w, reps)
}

func (api *API) ContractVestingEdit(w http.ResponseWriter, r *http.Request) {
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

	var data models.Vesting
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

	reps, err := service.ContractVestingEdit(user, contractAddress, data)
	if err != nil {
		//Unwrap apperror
		err, IsAppErr := apperrors.Unwrap(err)
		if !IsAppErr {
			log.Error("ContractVestingEdit error: ", zap.Error(err))
		}

		response.JsonError(w, err)
		return
	}

	response.Json(w, reps)
}

func (api *API) RemoveContractVesting(w http.ResponseWriter, r *http.Request) {
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

	var data models.Vesting
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

	err = service.RemoveContractVesting(user, contractAddress, data)
	if err != nil {
		//Unwrap apperror
		err, IsAppErr := apperrors.Unwrap(err)
		if !IsAppErr {
			log.Error("RemoveContractVesting error: ", zap.Error(err))
		}

		response.JsonError(w, err)
		return
	}

	response.Json(w, map[string]interface{}{"message": "success"})
}
