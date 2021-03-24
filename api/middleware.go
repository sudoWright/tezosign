package api

import (
	"context"
	"errors"
	"net/http"
	"tezosign/api/response"
	"tezosign/common/apperrors"
	"tezosign/common/log"
	"tezosign/repos"
	"tezosign/services"
	"tezosign/types"

	"go.uber.org/zap"

	"github.com/gorilla/mux"
)

type ContextKey string

const (
	ContextUserPubKey        ContextKey = "user_pubkey"
	ContextNetworkKey        ContextKey = "network"
	ContextNetworkContextKey ContextKey = "network_context"
)

func (api *API) RequireJWT(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	net, err := ToNetwork(mux.Vars(r)["network"])
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadParam, "network"))
		return
	}

	auth, err := api.provider.GetAuthProvider(net)
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadParam))
		return
	}

	userPubKey, err := auth.CheckSignatureAndGetUserPubKey(r)
	if err != nil {
		response.JsonError(w, err)
		return
	}

	if userPubKey == "" {
		response.JsonError(w, apperrors.New(apperrors.ErrService))
		return
	}

	typedPubKey := types.PubKey(userPubKey)

	err = typedPubKey.Validate()
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadParam))
		return
	}

	ctx := r.Context()
	ctx = context.WithValue(ctx, ContextUserPubKey, typedPubKey)
	r = r.WithContext(ctx)

	next(w, r)
}

func (api *API) CheckAndLoadNetwork(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	net, err := ToNetwork(mux.Vars(r)["network"])
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadParam, "network"))
		return
	}

	networkContext, err := api.provider.GetNetworkContext(net)
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadParam, err.Error()))
		return
	}

	ctx := r.Context()
	ctx = context.WithValue(ctx, ContextNetworkKey, net)
	ctx = context.WithValue(ctx, ContextNetworkContextKey, networkContext)
	r = r.WithContext(ctx)

	next(w, r)
}

const ContractIDParam = "contract_id"

func (api *API) OwnerAllowance(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	user, net, networkContext, err := GetUserNetworkContext(r)
	if err != nil {
		log.Error("Probably OwnerAllowance middleware was called without auth")
		response.JsonError(w, err)
		return
	}

	contractID := types.Address(mux.Vars(r)[ContractIDParam])
	if contractID == "" {
		log.Error("Probably OwnerAllowance middleware without contract_id param")
		response.JsonError(w, errors.New("Contract_id param not found"))
		return
	}

	if err = contractID.Validate(); err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadParam, ContractIDParam))
		return
	}

	service := services.New(repos.New(networkContext.Db), repos.New(networkContext.IndexerDB), networkContext.Client, networkContext.Auth, net)

	isOwner, err := service.GetUserAllowance(user, contractID)
	if err != nil {
		if err != nil {
			//Unwrap apperror
			err, IsAppErr := apperrors.Unwrap(err)
			if !IsAppErr {
				log.Error("GetUserAllowance error: ", zap.Error(err))
			}

			response.JsonError(w, err)
			return
		}
	}

	if !isOwner {
		response.JsonError(w, apperrors.NewWithDesc(apperrors.ErrNotAllowed, "pubkey not contains in storage"))
		return
	}

	next(w, r)
}
