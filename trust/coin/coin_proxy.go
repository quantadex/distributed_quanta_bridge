package coin

type CoinProxy struct {
	coin *DummyCoin
}

func (c *CoinProxy) Attach(coinName string) error {
	c.coin = GetDummyInstance()
	return nil
}

func (c *CoinProxy) GetTopBlockID() (int64, error) {
	return c.coin.GetTopBlockID()
}

func (c *CoinProxy) GetDepositsInBlock(blockID int64, trustAddress map[string]string) ([]*Deposit, error) {
	return c.coin.GetDepositsInBlock(blockID, trustAddress)
}

func (c *CoinProxy) SendWithdrawal(apiAddress string, w Withdrawal, s []byte) error {
	panic("implement me")
}


func (c *CoinProxy) EncodeRefund(w Withdrawal) (string, error) {
	return "", nil
}