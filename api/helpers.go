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
	}

	return "", fmt.Errorf("not supported network")
}
