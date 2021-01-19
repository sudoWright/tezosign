package services

import (
	"context"
	"msig/models"
	"msig/repos/contract"
)

type (
	Service interface {
	}

	// Provider is the abstract interface to get any repository.
	Provider interface {
		Health() error
		GetContract() contract.Repo
	}

	RPCProvider interface {
		ChainID(ctx context.Context) (chainID string, err error)
		Storage(ctx context.Context, contractAddress string) (storage string, err error)
		Script(context.Context, string) (models.BigMap, error)
		ManagerKey(ctx context.Context, address string) (pubKey string, err error)
	}

	ServiceFacade struct {
		repoProvider Provider
		rpcClient    RPCProvider
		net          models.Network
	}
)

func New(rp Provider, rpcClient RPCProvider, net models.Network) *ServiceFacade {

	return &ServiceFacade{
		repoProvider: rp,
		rpcClient:    rpcClient,
		net:          net,
	}
}
