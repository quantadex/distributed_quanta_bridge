package test

type QuantaNodeSecrets struct {
	NodeSecrets   []string
	SourceAccount string
}

type EthereumTrustSecrets struct {
	NodeSecrets   []string
	TrustContract string
}

type EthereumEnv struct {
	Rpc       string
	NetworkId string
}

type BtcSecrets struct {
	NodeSecrets []string
}

type LtcSecrets struct {
	NodeSecrets []string
}

var QUANTA_ISSUER = &QuantaNodeSecrets{
	NodeSecrets: []string{
		"ZBHK5VE5ZM5MJI3FM7JOW7MMUF3FIRUMV3BTLUTJWQHDFEN7MG3J4VAV",
		"ZDX6DGXBYAR3Z2BS4T4ITRTWPNJOSR5TPTVYN65UKEGP4ILOZ5GXU2KE",
		"ZC4U5P5DWNXGRUENOCOKZFHAWFKBE7JFOB2BCEKCM7BKXXKQE3DARXIJ",
	},
	SourceAccount: "QCISRUJ73RQBHB3C4LA6X537LPGSFZF3YUZ6MOPUOUJR5A63I5TLJML4",
}

var GRAPHENE_ISSUER = &QuantaNodeSecrets{
	NodeSecrets: []string{
		"5Jd9vxNwWXvMnBpcVm58gwXkJ4smzWDv9ChiBXwSRkvCTtekUrx",
		"5KFJnRn38wuXnpKGvkxmsyiWUuUkPXKZGvdG8aTzHCTvJMUQ4sA",
	},
	SourceAccount: "crosschain2",
}

var ROPSTEN_TRUST = &EthereumTrustSecrets{
	NodeSecrets: []string{
		// 0xba420ef5d725361d8fdc58cb1e4fa62eda9ec990
		"A7D7C6A92361590650AD0965970E186179F24F36B2B51CFE83F3AE8886BB6773",
		// 0xe0006458963c3773b051e767c5c63fee24cd7ff9
		"4C7F96D0CB8F2C48FD22CCB974513E6E9B0DC89475286BB24D2010E8D82AA461",
		// 0xba7573c0e805ef71acb7f1c4a55e7b0af416e96a
		"2E563A40747FA56419FB168ADF507C596E1A604D073D0F9E646B803DFA5BE94C",
	},
	TrustContract: "0xBD770336fF47A3B61D4f54cc0Fb541Ea7baAE92d",
}

var BTCSECRETS = &BtcSecrets{
	NodeSecrets: []string{
		"cNxQax7BfpbikeuCebPGCgTefTah5h1XhVDfaotVdFmXtaLCWLd9",
		"cUixT9PYjTtNzcVjF8sB7iM9JeEf8tLHm9Wjgo972x8opCRNTasS",
		"cPXngzEsUFpNCJ9DGYWyFLfCuGjzhsuM8N3sUf5z4HqLUUGuGp2h",
	},
}

var LTCSECRETS = &LtcSecrets{
	NodeSecrets: []string{
		"92P5DpWDiuttphtXV5qrHjMnFU2nAyiR8NpyEkF5s8uAngVgBFb",
		"926mkZAmMowq4HaLqpNjwuJuPe3vP6iTVQnt1x9GWdwbnwQjDea",
	},
}

var GRAPHENE_TRUST = &EthereumTrustSecrets{
	NodeSecrets: []string{
		// 0xba420ef5d725361d8fdc58cb1e4fa62eda9ec990
		"84d6b0af365017053af910682ebfccc36c34a1d5fff749471f1b532f86e144dd",
		// 0xe0006458963c3773b051e767c5c63fee24cd7ff9
		"5bebda860b34d4693f25af2da82332b2c89268e28566a9e0612c496002740d0c",
	},
	TrustContract: "0xBD770336fF47A3B61D4f54cc0Fb541Ea7baAE92d",
}

const ROPSTEN = "ROPSTEN"
const LOCAL = "LOCAL"

// must match up with the HorizonUrl
const QUANTA_ACCOUNT = "QCAO4HRMJDGFPUHRCLCSWARQTJXY2XTAFQUIRG2FAR3SCF26KQLAWZRN"

var ETHER_NETWORKS = map[string]EthereumEnv{
	ROPSTEN: EthereumEnv{"https://ropsten.infura.io/v3/7b880b2fb55c454985d1c1540f47cbf6", "3"},
	LOCAL:   EthereumEnv{"http://localhost:7545", "15"},
}
