package api

import (
	"msig/common/apperrors"
	"msig/infrustructure"
	"msig/models"
	"msig/types"
	"net/http"
)

//Get context values after middleware
func GetUserNetworkContext(r *http.Request) (user types.Address, net models.Network, networkContext infrustructure.NetworkContext, err error) {
	var ok bool
	ctx := r.Context()
	user, ok = ctx.Value(ContextUserKey).(types.Address)
	if !ok || (user.Validate() != nil) {
		return user, net, networkContext, apperrors.New(apperrors.ErrService)
	}

	net, networkContext, err = GetNetworkContext(r)
	if err != nil {
		return user, net, networkContext, err
	}

	return user, net, networkContext, nil
}

func GetNetworkContext(r *http.Request) (net models.Network, networkContext infrustructure.NetworkContext, err error) {
	var ok bool
	ctx := r.Context()
	net, ok = ctx.Value(ContextNetworkKey).(models.Network)
	if !ok {
		return net, networkContext, apperrors.New(apperrors.ErrService)
	}

	networkContext, ok = ctx.Value(ContextNetworkContextKey).(infrustructure.NetworkContext)
	if !ok {
		return net, networkContext, apperrors.New(apperrors.ErrService)
	}

	return net, networkContext, nil
}
