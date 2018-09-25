package coin

type CoinProxy struct {
	coin *DummyCoin
}

func (c *CoinProxy) Attach(coinName string) error {
	return nil
}

func (c *CoinProxy) GetTopBlockID() (int, error) {
	return c.coin.GetTopBlockID()
}

func (c *CoinProxy) GetDepositsInBlock(blockID int, trustAddress string) ([]*Deposit, error) {
	return c.coin.GetDepositsInBlock(blockID, trustAddress)
}

func (c *CoinProxy) SendWithdrawal(apiAddress string, w Withdrawal, s []byte) error {
	panic("implement me")
}


