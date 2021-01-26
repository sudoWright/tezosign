package infrustructure

import (
	"fmt"
	"msig/repos/postgres"
	"msig/services/auth"
	"msig/services/rpc_client"

	"gorm.io/gorm"
	"msig/conf"
	"msig/models"
)

type NetworkContext struct {
	Db     *gorm.DB
	Auth   *auth.Auth
	Client *rpc_client.Tezos
}

type Provider struct {
	networks map[models.Network]NetworkContext
}

func New(configs []conf.Network) (*Provider, error) {
	provider := &Provider{
		networks: make(map[models.Network]NetworkContext),
	}
	for i := range configs {

		db, err := postgres.New(configs[i].Params)
		if err != nil {
			return nil, err
		}

		rpcClient := rpc_client.New(configs[i].NodeRpc, configs[i].Name, configs[i].Name != models.NetworkMain)

		authProvider, err := auth.NewAuthProvider(configs[i].Auth, configs[i].Name)
		if err != nil {
			return nil, err
		}

		provider.networks[configs[i].Name] = NetworkContext{
			Db:     db,
			Auth:   authProvider,
			Client: rpcClient,
		}
	}
	return provider, nil
}

func (p *Provider) Close() {
	for _, v := range p.networks {
		sqlDB, err := v.Db.DB()
		if err != nil {
			return
		}
		sqlDB.Close()
	}
}

func (p *Provider) EnableTraceLevel() {
	for _, v := range p.networks {
		v.Db = v.Db.Debug()
	}
}

func (p *Provider) GetDb(net models.Network) (*gorm.DB, error) {
	if netcont, ok := p.networks[net]; ok {
		return netcont.Db, nil
	}
	return nil, fmt.Errorf("not enabled network")
}

func (p *Provider) GetRPCClient(net models.Network) (*rpc_client.Tezos, error) {
	if netcont, ok := p.networks[net]; ok {
		return netcont.Client, nil
	}
	return nil, fmt.Errorf("not enabled network")
}

func (p *Provider) GetAuthProvider(net models.Network) (*auth.Auth, error) {
	if netcont, ok := p.networks[net]; ok {
		return netcont.Auth, nil
	}
	return nil, fmt.Errorf("not enabled network")
}
