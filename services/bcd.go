package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"tezosign/models"
	"tezosign/types"
)

const betterCallDevAccountAPI = "https://better-call.dev/v1/account/edo2net/"

func getAccountTokensBalance(account types.Address) (balances models.AssetBalances, err error) {

	resp, err := http.Get(fmt.Sprint(betterCallDevAccountAPI, account.String()))
	if err != nil {
		return balances, err
	}

	if resp.StatusCode != http.StatusOK {
		return balances, fmt.Errorf("Not OK status code: %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return balances, fmt.Errorf("ReadAll error: %s", err.Error())
	}

	err = json.Unmarshal(body, &balances)
	if err != nil {
		return balances, fmt.Errorf("Unmarshal into baseResponse: %s", err.Error())
	}

	return balances, nil
}
