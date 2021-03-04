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

func (s *ServiceFacade) AssetsList(userPubKey types.PubKey, contractAddress types.Address) (assets []models.Asset, err error) {

	//TODO init limit from request
	limit := 100
	contract, isFound, err := s.repoProvider.GetContract().GetContract(contractAddress)
	if err != nil {
		return assets, err
	}

	if !isFound {
		return assets, apperrors.New(apperrors.ErrNotFound, "contract")
	}

	isOwner, err := s.GetUserAllowance(userPubKey, contractAddress)
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

	for i := range assets {
		if assets[i].ContractID.Valid {
			continue
		}

		assets[i].IsGlobal = true
	}

	return assets, nil
}

func (s *ServiceFacade) AssetsExchangeRates(userPubKey types.PubKey, contractAddress types.Address) (assetsRates map[string]interface{}, err error) {
	assets, err := s.AssetsList(userPubKey, contractAddress)
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

func (s *ServiceFacade) ContractAsset(userPubKey types.PubKey, contractAddress types.Address, reqAsset models.Asset) (asset models.Asset, err error) {

	contract, isFound, err := s.repoProvider.GetContract().GetContract(contractAddress)
	if err != nil {
		return asset, err
	}

	if !isFound {
		return asset, apperrors.New(apperrors.ErrNotFound, "contract")
	}

	isOwner, err := s.GetUserAllowance(userPubKey, contractAddress)
	if err != nil {
		return asset, err
	}

	if !isOwner {
		return asset, apperrors.New(apperrors.ErrNotAllowed)
	}

	//Сheck contract for FA
	isFAAsset, err := s.checkFAStandart(reqAsset.Address.String())
	if err != nil {
		return asset, err
	}

	if !isFAAsset {
		return asset, apperrors.New(apperrors.ErrBadParam, "not FA asset")
	}

	assetRepo := s.repoProvider.GetAsset()
	asset, isFound, err = assetRepo.GetAsset(contract.ID, reqAsset.Address)
	if err != nil {
		return asset, err
	}

	//Already created
	if isFound {
		return asset, apperrors.New(apperrors.ErrAlreadyExists, "asset")
	}

	reqAsset.ContractID = sql.NullInt64{
		Int64: int64(contract.ID),
		Valid: true,
	}

	err = assetRepo.CreateAsset(reqAsset)
	if err != nil {
		return asset, err
	}

	return reqAsset, nil
}

func (s *ServiceFacade) ContractAssetEdit(userPubKey types.PubKey, contractAddress types.Address, reqAsset models.Asset) (asset models.Asset, err error) {

	contract, isFound, err := s.repoProvider.GetContract().GetContract(contractAddress)
	if err != nil {
		return asset, err
	}

	if !isFound {
		return asset, apperrors.New(apperrors.ErrNotFound, "contract")
	}

	isOwner, err := s.GetUserAllowance(userPubKey, contractAddress)
	if err != nil {
		return asset, err
	}

	if !isOwner {
		return asset, apperrors.New(apperrors.ErrNotAllowed)
	}

	assetRepo := s.repoProvider.GetAsset()
	asset, isFound, err = assetRepo.GetAsset(contract.ID, reqAsset.Address)
	if err != nil {
		return asset, err
	}

	//Not created
	if !isFound {
		return asset, apperrors.New(apperrors.ErrNotFound, "asset")
	}

	//Global asset cannot be edited
	if !asset.ContractID.Valid {
		return asset, apperrors.New(apperrors.ErrNotAllowed, "global asset")
	}

	reqAsset.ContractID = sql.NullInt64{
		Int64: int64(contract.ID),
		Valid: true,
	}

	reqAsset.ID = asset.ID

	err = assetRepo.UpdateAsset(reqAsset)
	if err != nil {
		return asset, err
	}

	return reqAsset, nil
}

func (s *ServiceFacade) RemoveContractAsset(userPubKey types.PubKey, contractAddress types.Address, asset models.Asset) (err error) {

	contract, isFound, err := s.repoProvider.GetContract().GetContract(contractAddress)
	if err != nil {
		return err
	}

	if !isFound {
		return apperrors.New(apperrors.ErrNotFound, "contract")
	}

	isOwner, err := s.GetUserAllowance(userPubKey, contractAddress)
	if err != nil {
		return err
	}

	if !isOwner {
		return apperrors.New(apperrors.ErrNotAllowed)
	}

	assetRepo := s.repoProvider.GetAsset()
	asset, isFound, err = assetRepo.GetAsset(contract.ID, asset.Address)
	if err != nil {
		return err
	}

	if !isFound {
		return apperrors.New(apperrors.ErrNotFound, "asset")
	}

	//Global asset cannot be removed
	if !asset.ContractID.Valid {
		return apperrors.New(apperrors.ErrNotAllowed, "global asset")
	}

	err = assetRepo.DeleteContractAsset(asset.ID)
	if err != nil {
		return err
	}

	return nil
}
