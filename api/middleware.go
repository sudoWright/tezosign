package api

import (
	"context"
	"github.com/gorilla/mux"
	"net/http"
	"tezosign/api/response"
	"tezosign/common/apperrors"
	"tezosign/types"
)

type ContextKey string

const (
	ContextUserKey           ContextKey = "user_address"
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

	userAddress, err := auth.CheckSignatureAndGetUserAddress(r)
	if err != nil {
		response.JsonError(w, err)
		return
	}

	if userAddress == "" {
		response.JsonError(w, apperrors.New(apperrors.ErrService))
		return
	}

	typedAddress := types.Address(userAddress)

	err = typedAddress.Validate()
	if err != nil {
		response.JsonError(w, apperrors.New(apperrors.ErrBadParam))
		return
	}

	ctx := r.Context()
	ctx = context.WithValue(ctx, ContextUserKey, typedAddress)
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

func (api *API) OwnerAllowance(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	//TODO add user contract allowance
}
