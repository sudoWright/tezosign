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

	assets, err = s.repoProvider.GetAsset().GetAssetsList(contract.ID, isOwner, true, false)
	if err != nil {
		return assets, err
	}

	tokensMap, err := s.getContractTokensBalancesMap(contractAddress)
	if err != nil {
		return assets, err
	}

	for i := range assets {
		assets[i].Balances = tokensMap[assets[i].Address]

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

		eps, err := contract.InitAnnotsEntrypoints(script.Code.Storage)
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

	//Сheck contract for FA1.2 or FA2
	isFAAsset, err := s.checkFAStandart(reqAsset.Address, reqAsset.ContractType)
	if err != nil {
		return asset, err
	}

	if !isFAAsset {
		return asset, apperrors.New(apperrors.ErrBadParam, "not FA asset")
	}

	if reqAsset.ContractType == models.TypeFA12 && reqAsset.TokenID != 0 {
		return asset, apperrors.New(apperrors.ErrBadParam, "token_id")
	}

	assetRepo := s.repoProvider.GetAsset()
	asset, isFound, err = assetRepo.GetAsset(contract.ID, reqAsset.Address, reqAsset.TokenID)
	if err != nil {
		return asset, err
	}

	//Already created
	if isFound {
		//Check that asset is disabled
		if !asset.IsActive {
			err = assetRepo.EnableContractAsset(asset.ID)
			if err != nil {
				return asset, err
			}

			return asset, nil
		}

		return asset, apperrors.New(apperrors.ErrAlreadyExists, "asset")
	}

	reqAsset.ContractID = sql.NullInt64{
		Int64: int64(contract.ID),
		Valid: true,
	}

	tokensMap, err := s.getContractTokensBalancesMap(contractAddress)
	if err != nil {
		return asset, err
	}

	reqAsset.Balances = tokensMap[reqAsset.Address]

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

	assetRepo := s.repoProvider.GetAsset()
	asset, isFound, err = assetRepo.GetAsset(contract.ID, reqAsset.Address, reqAsset.TokenID)
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

	tokensMap, err := s.getContractTokensBalancesMap(contractAddress)
	if err != nil {
		return asset, err
	}

	reqAsset.Balances = tokensMap[reqAsset.Address]

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

	assetRepo := s.repoProvider.GetAsset()
	asset, isFound, err = assetRepo.GetAsset(contract.ID, asset.Address, asset.TokenID)
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

	err = assetRepo.DisableContractAsset(asset.ID)
	if err != nil {
		return err
	}

	return nil
}

const transferEntrypoint = "transfer"

func (s *ServiceFacade) AssetsIncomeOperations() (count uint64, err error) {

	//Get assets
	assets, err := s.repoProvider.GetAsset().GetAssetsList(0, false, true, true)
	if err != nil {
		return count, err
	}

	//TODO add limit
	limit := 1000
	contracts, err := s.repoProvider.GetContract().GetContractsList(limit, 0)
	if err != nil {
		return count, err
	}

	contractsMap := make(map[types.Address]models.Contract, len(contracts))

	for i := range contracts {
		contractsMap[contracts[i].Address] = contracts[i]
	}

	networkID, err := s.rpcClient.ChainID(context.Background())
	if err != nil {
		return count, err
	}

	for i := range assets {

		operationsCount, err := s.processAssetOperations(contractsMap, networkID, assets[i])
		if err != nil {
			return count, err
		}

		count += operationsCount
	}

	return count, err
}

func (s *ServiceFacade) GetAssetMetadata(assetID types.Address, tokenID int64) (resp interface{}, err error) {

	indexerRepo := s.indexerRepoProvider.GetIndexer()

	script, isFound, err := indexerRepo.GetContractScript(assetID)
	if err != nil {
		return resp, err
	}

	if !isFound {
		return resp, apperrors.New(apperrors.ErrNotFound, "script")
	}

	e, err := contract.InitAnnotsEntrypoints(script.StorageSchema.MichelinePrim())
	if err != nil {
		return resp, err
	}

	entrypoint, ok := e[contract.MetaDataEntrypoint]
	if !ok {
		return resp, apperrors.New(apperrors.ErrBadRequest, "token_metadata")
	}

	storage, isFound, err := indexerRepo.GetContractStorage(assetID)
	if err != nil {
		return resp, err
	}

	if !isFound {
		return resp, apperrors.New(apperrors.ErrNotFound, "storage")
	}

	value, err := contract.GetStorageValue(entrypoint, storage.RawValue.MichelinePrim())
	if err != nil {
		return resp, err
	}

	bigMap := value.Int.Int64()

	hash, err := contract.GetBigMapKeyHash(tokenID)
	if err != nil {
		return resp, err
	}

	bigMapValue, isFound, err := s.rpcClient.BigMapKey(context.Background(), bigMap, hash)
	if err != nil {
		return resp, err
	}

	if !isFound {
		return resp, apperrors.New(apperrors.ErrNotFound, "big_map value")
	}

	v, err := contract.ParseMetadata(bigMapValue)
	if err != nil {
		return resp, err
	}

	return v, nil
}

func (s *ServiceFacade) processAssetOperations(contractsMap map[types.Address]models.Contract, networkID string, asset models.Asset) (count uint64, err error) {

	transferType := models.IncomeFATransfer
	if asset.ContractType == models.TypeFA2 {
		transferType = models.IncomeFA2Transfer
	}

	assetOperations, err := s.indexerRepoProvider.GetIndexer().GetContractOperations(asset.Address, asset.LastOperationBlockLevel, transferEntrypoint)
	if err != nil {
		return count, err
	}

	//New operations not founds
	//Setup last block level as last known block
	if len(assetOperations) == 0 {
		block, err := s.indexerRepoProvider.GetIndexer().GetLastBlock()
		if err != nil {
			return count, err
		}

		asset.LastOperationBlockLevel = block.Level
	}

	s.repoProvider.Start(context.Background())
	defer s.repoProvider.RollbackUnlessCommitted()

	for j := range assetOperations {
		txs := contract.AssetOperation(assetOperations[j].RawParameters.MichelinePrim(), asset.ContractType)

		transferUnits := groupOperations(asset.TokenID, contractsMap, txs)
		for contractAddress, transfers := range transferUnits {

			err = s.repoProvider.GetContract().SavePayload(models.Request{
				Hash:       operationID(assetOperations[j].OpHash),
				ContractID: contractsMap[contractAddress].ID,
				Counter:    nil,
				Status:     models.StatusSuccess,
				CreatedAt:  assetOperations[j].Timestamp,
				Info: models.ContractOperationRequest{
					ContractID:   contractAddress,
					AssetID:      asset.Address,
					Type:         transferType,
					TransferList: transfers,
				},
				NetworkID:   networkID,
				OperationID: &assetOperations[j].OpHash,
			})

			//Increment counter of saved operations
			count++
		}

		asset.LastOperationBlockLevel = assetOperations[j].Level
	}

	err = s.repoProvider.GetAsset().UpdateAsset(asset)
	if err != nil {
		return 0, err
	}

	err = s.repoProvider.Commit()
	if err != nil {
		return 0, err
	}

	return count, nil
}

func groupOperations(tokenID uint64, contractsMap map[types.Address]models.Contract, txs []models.TransferUnit) map[types.Address][]models.TransferUnit {
	//map Address TO map address From
	transferUnitsGroupsByFromAddress := map[types.Address]map[types.Address]models.TransferUnit{}

	for k := range txs {
		from := txs[k].From

		for _, value := range txs[k].Txs {

			//Skip txs to not our contracts
			_, ok := contractsMap[value.To]
			if !ok {
				continue
			}

			//Init internal map
			if _, ok = transferUnitsGroupsByFromAddress[value.To]; !ok {
				transferUnitsGroupsByFromAddress[value.To] = map[types.Address]models.TransferUnit{}
			}

			//Skip txs with another txID
			if value.TokenID != tokenID {
				continue
			}
			//Init first value
			_, ok = transferUnitsGroupsByFromAddress[value.To][from]
			if !ok {
				transferUnitsGroupsByFromAddress[value.To][from] = models.TransferUnit{
					From: from,
					Txs:  []models.Tx{value},
				}
				continue
			}

			//Append tx
			transferUnitsGroupsByFromAddress[value.To][from] = models.TransferUnit{
				From: from,
				Txs:  append(transferUnitsGroupsByFromAddress[value.To][from].Txs, value),
			}
		}
	}

	transferUnits := map[types.Address][]models.TransferUnit{}

	//Merge to map Address TO
	for address, mapFrom := range transferUnitsGroupsByFromAddress {

		units := make([]models.TransferUnit, 0, len(mapFrom))

		for _, value := range mapFrom {
			units = append(units, value)
		}

		transferUnits[address] = append(transferUnits[address], units...)
	}

	return transferUnits
}

func (s *ServiceFacade) getContractTokensBalancesMap(contractAddress types.Address) (tokensMap map[types.Address][]models.TokenBalance, err error) {

	balances, err := getAccountTokensBalance(contractAddress, s.net)
	if err != nil {
		return tokensMap, err
	}

	//TODO process total

	tokensMap = make(map[types.Address][]models.TokenBalance, len(balances.Balances))
	for i := range balances.Balances {
		tokensMap[balances.Balances[i].Asset] = append(tokensMap[balances.Balances[i].Asset], balances.Balances[i].TokenBalance)
	}

	return tokensMap, nil
}
