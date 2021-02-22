package services

import (
	"context"
	"tezosign/models"
	"tezosign/services/contract"

	"github.com/wedancedalot/decimal"
)

const (
	TezosPrecision    = 6
	TruncatePrecision = 8
)

func (s *ServiceFacade) AssetsList() (assets []models.Asset, er error) {

	//TODO init limit from request
	limit := 100

	assets, err := s.repoProvider.GetAsset().GetAssetsList(limit, 0)
	if err != nil {
		return assets, err
	}

	return assets, nil
}

func (s *ServiceFacade) AssetsExchangeRates() (assetsRates map[string]interface{}, er error) {
	//TODO init limit from request
	limit := 100

	assets, err := s.repoProvider.GetAsset().GetAssetsList(limit, 0)
	if err != nil {
		return assetsRates, err
	}

	//Init map
	assetsRates = make(map[string]interface{}, len(assets))

	for i := range assets {
		//Skip assets not presented on Exchange
		if len(assets[i].DexterAddress) == 0 {
			continue
		}

		script, err := s.rpcClient.Script(context.Background(), assets[i].DexterAddress)
		if err != nil {
			return assetsRates, err
		}

		eps, err := contract.InitStorageAnnotsEntrypoints(script.Code.Storage)
		if err != nil {
			return assetsRates, err
		}

		tokenPool, err := contract.GetDexterContractTokenPool(eps, script.Storage)
		if err != nil {
			return assetsRates, err
		}

		//In mutez
		xtzPool, err := contract.GetDexterContractXTZPool(eps, script.Storage)
		if err != nil {
			return assetsRates, err
		}

		tPool := decimal.NewFromBigInt(tokenPool, -int32(assets[i].Scale))

		xPool := decimal.NewFromBigInt(xtzPool, -TezosPrecision)

		price := tPool.Div(xPool)

		assetsRates[assets[i].Ticker] = price.Truncate(TruncatePrecision)
	}

	return
}
