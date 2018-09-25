package quanta

import "github.com/quantadex/distributed_quanta_bridge/trust/peer_contact"

type QuantaClient struct {

}

func (q *QuantaClient) Attach() error {
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


