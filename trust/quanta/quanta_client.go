package quanta

import (
	"github.com/quantadex/distributed_quanta_bridge/trust/peer_contact"
	"github.com/quantadex/distributed_quanta_bridge/trust/coin"
	"github.com/spf13/viper"
	b "github.com/quantadex/stellar_go/build"
	"github.com/quantadex/stellar_go/clients/horizon"

	"net/http"
	"fmt"
)

type QuantaClient struct {
	network string
	issuer string // pub key
	horizonClient *horizon.Client
}

// remember to test coins < 10^7
func (q *QuantaClient) CreateProposeTransaction(deposit *coin.Deposit) (string, error) {
	amount := fmt.Sprintf("%.7f",float64(deposit.Amount)/10000000)

	tx, err := b.Transaction(
		b.Network{q.network},
		b.SourceAccount{q.issuer},
		//b.AutoSequence{horizon.DefaultTestNetClient},
		b.Sequence{ 0 },
		b.Payment(
			b.Destination{deposit.QuantaAddr},
			b.NativeAmount{ amount},
		),
	)

	if err != nil {
		return "", err
	}

	txe, err := tx.Sign()

	if err != nil {
		return "", err
	}

	return txe.Base64()
}

func (q *QuantaClient) Attach() error {
	q.network = viper.GetString("NETWORK_PASSPHRASE")
	viper.GetString("NETWORK_PASSPHRASE")
	q.issuer = viper.GetString("ISSUER_ADDRESS")
	q.horizonClient = &horizon.Client{
		URL:  viper.GetString("HORIZON_URL"),
		HTTP: http.DefaultClient,
	}
	return nil
}

func (q *QuantaClient) GetTopBlockID() (int, error) {
	return 0, nil
}

func (q *QuantaClient) GetRefundsInBlock(blockID int, trustAddress string) ([]Refund, error) {
	return nil, nil
}

func (q *QuantaClient) ProcessDeposit(deposit peer_contact.PeerMessage) error {
	return nil
}




