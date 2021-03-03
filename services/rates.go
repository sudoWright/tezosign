package services

import "tezosign/models"

func (s *ServiceFacade) TezosExchangeRates() (quote models.Quote, err error) {
	quote, err = s.indexerRepoProvider.GetIndexer().GetTezosQuote()
	if err != nil {
		return quote, err
	}

	return quote, nil
}
