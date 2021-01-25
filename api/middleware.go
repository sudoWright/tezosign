package api

import (
	"context"
	"github.com/gorilla/mux"
	"msig/api/response"
	"msig/common/apperrors"
	"msig/types"
	"net/http"
)

const ContextUserKey = "user_address"

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
