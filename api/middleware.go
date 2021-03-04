package api

import (
	"context"
	"net/http"
	"tezosign/api/response"
	"tezosign/common/apperrors"
	"tezosign/types"

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

func (api *API) OwnerAllowance(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	//TODO add user contract allowance
}
