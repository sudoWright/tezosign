package services

import (
	"context"
	"tezosign/models"
	"tezosign/repos/auth"
	"tezosign/repos/contract"
	"tezosign/repos/indexer"
	"tezosign/types"

	"blockwatch.cc/tzindex/micheline"
)

type (
	Service interface {
	}

	//Db transaction
	DBTx interface {
		Start(ctx context.Context)
		RollbackUnlessCommitted()
		Commit() error
	}

	// Provider is the abstract interface to get any repository.
	Provider interface {
		Health() error
		GetContract() contract.Repo
		GetAuth() auth.Repo

		DBTx
	}

	IndexerProvider interface {
		GetIndexer() indexer.Repo
	}

	RPCProvider interface {
		ChainID(ctx context.Context) (chainID string, err error)
		Storage(ctx context.Context, contractAddress string) (storage string, err error)
		Script(context.Context, string) (micheline.Script, error)
		ManagerKey(ctx context.Context, address string) (pubKey string, err error)
		Balance(ctx context.Context, address string) (balance int64, err error)
	}

	AuthProvider interface {
		GenerateAuthTokens(address types.Address) (string, string, error)

		EncodeSessionCookie(data map[string]string) (string, error)
		DecodeSessionCookie(cookie string) (map[string]string, error)
	}

	ServiceFacade struct {
		repoProvider        Provider
		indexerRepoProvider IndexerProvider
		rpcClient           RPCProvider
		auth                AuthProvider
		net                 models.Network
	}
)

func New(rp Provider, iRp IndexerProvider, rpcClient RPCProvider, auth AuthProvider, net models.Network) *ServiceFacade {

	return &ServiceFacade{
		repoProvider:        rp,
		indexerRepoProvider: iRp,
		auth:                auth,
		rpcClient:           rpcClient,
		net:                 net,
	}
}
