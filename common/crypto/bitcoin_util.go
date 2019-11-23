package crypto

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcjson"
	"strings"
)

type BitcoinRPCClient interface {
	RawRequest(method string, params []json.RawMessage) (json.RawMessage, error)
}

func ValidateNetwork(client BitcoinRPCClient, expected string) error {
	res, err := client.RawRequest("getnetworkinfo", nil)
	if err != nil {
		return err
	}

	// Unmarshal result as a gettransaction result object
	var getTx btcjson.GetNetworkInfoResult
	err = json.Unmarshal(res, &getTx)
	if err != nil {
		return err
	}
	if !strings.Contains(getTx.SubVersion, expected) {
		return errors.New(fmt.Sprintf("Wrong blockchain, expecting %s but got %s", expected, getTx.SubVersion))
	}

	return nil
}
