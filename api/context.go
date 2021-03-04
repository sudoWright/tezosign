package api

import (
	"net/http"
	"tezosign/common/apperrors"
	"tezosign/infrustructure"
	"tezosign/models"
	"tezosign/types"
)

//Get context values after middleware
func GetUserNetworkContext(r *http.Request) (userPubKey types.PubKey, net models.Network, networkContext infrustructure.NetworkContext, err error) {
	var ok bool
	ctx := r.Context()
	userPubKey, ok = ctx.Value(ContextUserPubKey).(types.PubKey)
	if !ok || (userPubKey.Validate() != nil) {
		return userPubKey, net, networkContext, apperrors.New(apperrors.ErrService)
	}

	net, networkContext, err = GetNetworkContext(r)
	if err != nil {
		return userPubKey, net, networkContext, err
	}

	return userPubKey, net, networkContext, nil
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
