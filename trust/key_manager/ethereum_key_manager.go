package key_manager

type EthereumKeyManager struct{}

func (e *EthereumKeyManager) CreateNodeKeys() error {
	panic("implement me")
}

func (e *EthereumKeyManager) LoadNodeKeys(privKey string) error {
	panic("implement me")
}

func (e *EthereumKeyManager) GetPublicKey() (string, error) {
	panic("implement me")
}

func (e *EthereumKeyManager) SignMessage(original []byte) ([]byte, error) {
	panic("implement me")
}

func (e *EthereumKeyManager) SignMessageObj(original interface{}) (*string) {
	panic("implement me")
}

func (e *EthereumKeyManager) VerifySignatureObj(original interface{}, key string) bool {
	panic("implement me")
}

func (e *EthereumKeyManager) SignTransaction(base64 string) (string, error) {
	panic("implement me")
}

func (e *EthereumKeyManager) VerifyTransaction(base64 string) (bool, error) {
	panic("implement me")
}
