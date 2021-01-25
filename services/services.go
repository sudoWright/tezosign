package services

import (
	"context"
	"msig/models"
	"msig/repos/auth"
	"msig/repos/contract"
	"msig/types"
)

type (
	Service interface {
	}

	// Provider is the abstract interface to get any repository.
	Provider interface {
		Health() error
		GetContract() contract.Repo
		GetAuth() auth.Repo
	}

	RPCProvider interface {
		ChainID(ctx context.Context) (chainID string, err error)
		Storage(ctx context.Context, contractAddress string) (storage string, err error)
		Script(context.Context, string) (models.BigMap, error)
		ManagerKey(ctx context.Context, address string) (pubKey string, err error)
	}

	AuthProvider interface {
		GenerateAuthTokens(address types.Address) (string, string, error)

		EncodeSessionCookie(data map[string]string) (string, error)
		DecodeSessionCookie(cookie string) (map[string]string, error)
	}

	ServiceFacade struct {
		repoProvider Provider
		rpcClient    RPCProvider
		auth         AuthProvider
		net          models.Network
	}
)

func New(rp Provider, rpcClient RPCProvider, auth AuthProvider, net models.Network) *ServiceFacade {

	return &ServiceFacade{
		repoProvider: rp,
		auth:         auth,
		rpcClient:    rpcClient,
		net:          net,
	}
}
