package common

type Config struct {
	ExternalListenPort int
	ListenIp           string
	ListenPort         int
	UsePrevKeys        bool
	KvDbName           string
	DatabaseUrl        string
	CoinName           string
	IssuerAddress      string
	NodeKey            string
	NetworkUrl         string
	ChainId            string
	RegistrarIp        string
	RegistrarPort      int
	EthereumNetworkId  string
	EthereumBlockStart int64
	EthereumRpc        string
	EthereumTrustAddr  string
	EthereumKeyStore   string
	BtcPrivateKey	   string
	BtcBlockStart 		int64
	LogLevel           string
	MinNodes           int
}
