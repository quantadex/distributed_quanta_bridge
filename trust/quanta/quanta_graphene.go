package quanta

/**
 Links:
http://docs.bitshares.org/integration/traders/index.html#public-api
https://github.com/scorum/bitshares-go/blob/master/apis/database/api_test.go

 */
import (
	"github.com/quantadex/distributed_quanta_bridge/common/kv_store"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
)

type QuantaGraphene struct{}

func (q *QuantaGraphene) Attach() error {
	panic("implement me")
}

func (q *QuantaGraphene) AttachQueue(kv kv_store.KVStore) error {
	panic("implement me")
}

// get_dynamics
func (q *QuantaGraphene) GetTopBlockID(accountId string) (int64, error) {
	panic("implement me")
}

// get block , transfer
func (q *QuantaGraphene) GetRefundsInBlock(blockID int64, trustAddress string) ([]Refund, int64, error) {
	panic("implement me")
}

//submission
func (q *QuantaGraphene) ProcessDeposit(deposit *coin.Deposit, proposed string) error {
	panic("implement me")
}

func (q *QuantaGraphene) GetBalance(assetName string, quantaAddress string) (float64, error) {
	panic("implement me")
}

func (q *QuantaGraphene) GetAllBalances(quantaAddress string) (map[string]float64, error) {
	panic("implement me")
}

// https://github.com/scorum/bitshares-go/blob/bbfc9bedaa1b2ddaead3eafe47237efcd9b8496d/client.go
func (q *QuantaGraphene) CreateProposeTransaction(*coin.Deposit) (string, error) {
	panic("implement me")
}

func (q *QuantaGraphene) DecodeTransaction(base64 string) (*coin.Deposit, error) {
	panic("implement me")
}

