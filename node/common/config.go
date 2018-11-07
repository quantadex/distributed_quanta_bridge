package common

type Config struct {
	ExternalListenPort int
	ListenIp           string
	ListenPort         int
	UsePrevKeys        bool
	KvDbName           string
	CoinName           string
	IssuerAddress      string
	NodeKey            string
	HorizonUrl         string
	NetworkPassphrase  string
	RegistrarIp        string
	RegistrarPort      int
	EthereumNetworkId  string
	EthereumBlockStart int64
	EthereumRpc        string
	EthereumTrustAddr  string
	EthereumKeyStore   string
}