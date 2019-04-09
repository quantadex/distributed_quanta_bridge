package sync

import (
	"encoding/json"
	"fmt"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"math/big"
)

type LitecoinSync struct {
	DepositSync
	issuingSymbol map[string]string
}

func (c *LitecoinSync) Setup() {
	c.fnDepositInBlock = c.GetDepositsInBlock
	c.fnGetWatchAddress = c.GetWatchAddress
	c.fnTransformCoin = c.TransformCoin
}

func (c *LitecoinSync) TransformCoin(dep *coin.Deposit) *coin.Deposit {
	if dep.CoinName == "LTC" {
		dep.CoinName = c.issuingSymbol["ltc"]
	}
	dep.Amount = coin.PowerDelta(*big.NewInt(dep.Amount), 8, int(c.coinInfo[c.issuingSymbol["ltc"]].Precision))
	return dep
}

func (c *LitecoinSync) GetDepositsInBlock(blockID int64) ([]*coin.Deposit, error) {
	watchAddresses := c.rDb.GetCrosschainByBlockchain(coin.BLOCKCHAIN_LTC)
	watchMap := make(map[string]string)

	for _, w := range watchAddresses {
		watchMap[w.Address] = w.QuantaAddr
	}
	deposits, err := c.coinChannel.GetDepositsInBlock(blockID, watchMap)

	for _, dep := range deposits {
		dep = c.TransformCoin(dep)
	}

	if err != nil {
		c.logger.Info("getDepositsInBlock failed " + err.Error())
		return nil, err
	}

	if len(deposits) > 0 {
		msg, _ := json.Marshal(deposits)
		fmt.Printf("events = %v\n", string(msg))
	}

	return deposits, err
}

func (c *LitecoinSync) GetWatchAddress() map[string]string {
	watchAddresses := c.rDb.GetCrosschainByBlockchain(coin.BLOCKCHAIN_LTC)
	watchMap := make(map[string]string)

	for _, w := range watchAddresses {
		watchMap[w.Address] = w.QuantaAddr
	}
	return watchMap
}

func (c *LitecoinSync) PostProcessBlock(blockID int64) error {
	panic("not imp")
}
