package services

import (
	"context"
	"database/sql"
	"tezosign/common/apperrors"
	"tezosign/models"
	"tezosign/services/contract"
	"tezosign/types"

	"github.com/wedancedalot/decimal"
)

const (
	TezosPrecision    = 6
	TruncatePrecision = 8
)

func (s *ServiceFacade) AssetsList(user, contractAddress types.Address) (assets []models.Asset, err error) {

	//TODO init limit from request
	limit := 100
	contract, isFound, err := s.repoProvider.GetContract().GetContract(contractAddress)
	if err != nil {
		return assets, err
	}

	if !isFound {
		return assets, apperrors.New(apperrors.ErrNotFound, "contract")
	}

	isOwner, err := s.GetUserAllowance(user, contractAddress)
	if err != nil {
		return assets, err
	}

	if !isOwner {
		return assets, apperrors.New(apperrors.ErrNotAllowed)
	}

	assets, err = s.repoProvider.GetAsset().GetAssetsList(contract.ID, limit, 0)
	if err != nil {
		return assets, err
	}

	return assets, nil
}

func (s *ServiceFacade) AssetsExchangeRates(user, contractAddress types.Address) (assetsRates map[string]interface{}, err error) {
	assets, err := s.AssetsList(user, contractAddress)
	if err != nil {
		return assetsRates, err
	}

	//Init map
	assetsRates = make(map[string]interface{}, len(assets))

	for i := range assets {
		//Skip assets not presented on Exchange
		if assets[i].DexterAddress == nil {
			continue
		}

		script, err := s.rpcClient.Script(context.Background(), *assets[i].DexterAddress)
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

func (s *ServiceFacade) ContractAsset(user, contractAddress types.Address, asset models.Asset) (models.Asset, error) {

	contract, isFound, err := s.repoProvider.GetContract().GetContract(contractAddress)
	if err != nil {
		return asset, err
	}

	if !isFound {
		return asset, apperrors.New(apperrors.ErrNotFound, "contract")
	}

	isOwner, err := s.GetUserAllowance(user, contractAddress)
	if err != nil {
		return asset, err
	}

	if !isOwner {
		return asset, apperrors.New(apperrors.ErrNotAllowed)
	}

	asset.ContractID = sql.NullInt64{
		Int64: int64(contract.ID),
		Valid: true,
	}

	//Сheck asset already not added
	assetRepo := s.repoProvider.GetAsset()
	_, isFound, err = assetRepo.GetAsset(contract.ID, types.Address(asset.Address))
	if err != nil {
		return asset, err
	}

	if isFound {
		return asset, apperrors.New(apperrors.ErrAlreadyExists, "asset")
	}

	//Сheck contract for FA
	isFAAsset, err := s.checkFAStandart(asset.Address)
	if err != nil {
		return asset, err
	}

	if !isFAAsset {
		return asset, apperrors.New(apperrors.ErrBadParam, "not FA asset")
	}

	//Skip dexter address from user req
	asset.DexterAddress = nil

	err = assetRepo.CreateAsset(asset)
	if err != nil {
		return asset, err
	}

	return asset, nil
}
