package api

import (
	"fmt"
	"strings"
	"tezosign/models"
)

func ToNetwork(net string) (models.Network, error) {
	switch strings.ToLower(net) {
	case "main", "mainnet":
		return models.NetworkMain, nil
	case "delphi", "delphinet":
		return models.NetworkDelphi, nil
	case "edo", "edonet":
		return models.NetworkEdo, nil
	}

	return "", fmt.Errorf("not supported network")
}
