package common

type Config struct {
	ExternalListenPort int
	ListenIp           string
	ListenPort         int
	UsePrevKeys        bool
	KvDbName           string

	CoinMapping          map[string]string
	IssuerAddress        string
	NetworkUrl           string
	ChainId              string
	RegistrarIp          string
	RegistrarPort        int
	EthereumNetworkId    string
	EthereumBlockStart   int64
	EthereumRpc          string
	EthereumTrustAddr    string
	EthMinConfirmation   int64
	EthDegradedThreshold int64
	EthFailureThreshold  int64

	BtcRpc               string
	BtcBlockStart        int64
	BtcNetwork           string
	BtcMinConfirmation   int64
	BtcDegradedThreshold int64
	BtcFailureThreshold  int64

	LogLevel      string
	MinNodes      int
	EthFlush      bool
	Erc20Mapping  map[string]string
	MinBlockReuse int64

	LtcRpc               string
	LtcNetwork           string
	LtcBlockStart        int64
	LtcMinConfirmation   int64
	LtcDegradedThreshold int64
	LtcFailureThreshold  int64

	BchRpc               string
	BchNetwork           string
	BchBlockStart        int64
	BchMinConfirmation   int64
	BchDegradedThreshold int64
	BchFailureThreshold  int64

	QuantaDegradedThreshold   int64
	QuantaFailureThreshold    int64
	DepDegradedThreshold      int64
	DepFailureThreshold       int64
	WithdrawDegradedThreshold int64
	WithdrawFailureThreshold  int64
}
