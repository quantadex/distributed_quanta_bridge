package common

type Secrets struct {
	NodeKey          string
	EthereumKeyStore string
	DatabaseUrl      string

	BtcPrivateKey  string
	BtcRpcUser     string
	BtcRpcPassword string
	BtcSigners     []string

	LtcPrivateKey  string
	LtcRpcUser     string
	LtcRpcPassword string
	LtcSigners     []string

	BchPrivateKey  string
	BchRpcUser     string
	BchRpcPassword string
	BchSigners     []string

	GrapheneSeedPrefix string
}
